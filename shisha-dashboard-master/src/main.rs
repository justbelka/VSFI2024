use std::env;

use axum::extract::State;
use axum::http::{header, StatusCode, Uri};
use axum::response::{Html, IntoResponse, Response};
use axum::Router;
use axum::routing::get;
use axum_prometheus::PrometheusMetricLayer;
use chrono::{DateTime, Local};
use rust_embed::RustEmbed;
use sea_orm::*;
use sea_orm_migration::prelude::*;
use serde::Serialize;
use tera::{Context, Tera};

use crate::entities::event::{EventType, Model};
use crate::messages::*;
use crate::migrator::Migrator;
use crate::service::query::Query as QueryCore;

mod migrator;
mod entities;
mod service;
mod messages;

#[derive(Clone)]
struct AppState {
    db_connection: DatabaseConnection,
    templates: Tera,
}

#[derive(Serialize)]
struct FrontendEvent {
    actor: String,
    event_type: EventType,
    event_date: DateTime<Local>,
}

impl Into<FrontendEvent> for &Model {
    fn into(self) -> FrontendEvent {
        FrontendEvent {
            actor: self.actor.clone(),
            event_type: self.event_type.clone(),
            event_date: DateTime::from(self.event_date.and_utc()),
        }
    }
}

#[tokio::main]
async fn main() -> anyhow::Result<()> {
    env_logger::builder()
        .filter_level(log::LevelFilter::Debug)
        .is_test(true)
        .init();
    dotenvy::dotenv().ok();

    let db_url = env::var("DATABASE_URL").expect("DATABASE_URL env var is not set");
    let boostrap_servers = env::var("BOOTSTRAP_SERVERS").expect("BOOTSTRAP_SERVERS var is not set");
    let topic_name = env::var("EVENTS_TOPIC").expect("EVENTS_TOPIC var is not set");
    let listener_config = ListenerConfig { boostrap_servers, topic: topic_name };
    let db_connection = Database::connect(db_url)
        .await
        .expect("Failed to connect to database");
    Migrator::up(&db_connection, None).await.expect("Failed to run migrations");
    let consumer = create_kafka_consumer(&listener_config);
    let connection_clone = db_connection.clone();
    tokio::spawn(async move {
        consume_events_task(consumer, &listener_config, &connection_clone).await
    });
    let mut tera = Tera::default();
    Templates::iter().for_each(|item| {
        let template_bytes = Templates::get(item.as_ref()).expect("Template Not Found").data;
        let template_content = std::str::from_utf8(template_bytes.as_ref()).expect("Failed to parse template content");
        tera.add_raw_template(item.as_ref(), template_content.as_ref()).expect("Failed to add template");
    });
    let (prometheus_layer, metric_handle) = PrometheusMetricLayer::pair();
    let state = AppState { db_connection, templates: tera };
    let app = Router::new()
        .route("/", get(analytics))
        .route("/assets/*path", get(serve_asset))
        .layer(prometheus_layer)
        .route("/metrics/prometheus", get(|| async move { metric_handle.render() }))
        .route("/health", get(|state: State<AppState>| async move {
            state.db_connection.ping().await.map(|_| { StatusCode::OK }).map_err(|e| {
                log::error!("Failed to ping database {}", e);
                StatusCode::INTERNAL_SERVER_ERROR
            })
        }))
        .route("/ready", get(|| async move { StatusCode::OK }))
        .with_state(state);

    let listener = tokio::net::TcpListener::bind("0.0.0.0:5050").await.expect("Failed to start Axum server!");
    axum::serve(listener, app).await?;
    Ok(())
}

async fn analytics(state: State<AppState>) -> Result<Html<String>, StatusCode> {
    let mut context = Context::new();
    let top_buyer = QueryCore::top_buyer_name(&state.db_connection).await.expect("DataError");
    let top_uploader = QueryCore::top_uploader(&state.db_connection).await.expect("DataError");
    let last_10_events: Vec<FrontendEvent> = QueryCore::last_10_events(&state.db_connection).await.expect("DataError")
        .iter()
        .map(|item| { item.into() })
        .collect();
    let money_earned = QueryCore::count_money(&state.db_connection).await.expect("DataError");
    let shisha_uploads = QueryCore::count_uploads(&state.db_connection).await.expect("DataError");

    context.insert("top_buyer", &top_buyer);
    context.insert("top_uploader", &top_uploader);
    context.insert("last_10_events", &last_10_events);
    context.insert("money_earned", &money_earned);
    context.insert("shisha_uploads", &shisha_uploads);

    Ok(Html(state.templates.render("index.html", &context).expect("template not rendered")))
}

async fn serve_asset(uri: Uri) -> impl IntoResponse {
    let mut path = uri.path().trim_start_matches('/').to_string();
    if path.starts_with("assets/") {
        path = path.replace("assets/", "")
    }
    StaticFile(path)
}

#[derive(RustEmbed)]
#[folder = "assets/"]
struct Asset;

#[derive(RustEmbed)]
#[folder = "templates/"]
struct Templates;

pub struct StaticFile<T>(pub T);

impl<T> IntoResponse for StaticFile<T>
where
    T: Into<String>,
{
    fn into_response(self) -> Response {
        let path = self.0.into();
        match Asset::get(path.as_str()) {
            Some(file) => {
                let mime = mime_guess::from_path(path).first_or_octet_stream();
                ([(header::CONTENT_TYPE, mime.as_ref())], file.data).into_response()
            }
            None => (StatusCode::NOT_FOUND, "404 Not Found").into_response()
        }
    }
}
