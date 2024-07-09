use std::time::Duration;
use chrono::{TimeDelta};
use sea_orm::Database;
use sea_orm_migration::prelude::*;
use shisha::entities::event::EventType;
use shisha::messages::ShishaMessage;
use shisha::service::query::Query as QueryCore;

use shisha::migrator::Migrator;

use crate::common::{db_url, listener_config, setup, write_shisha_message};

mod common;

#[tokio::test]
async fn can_write_message_to_db() {
    let (rp, mysql, producer) = setup().await;
    let db_url = db_url(&mysql).await;
    let db_connection = Database::connect(db_url)
        .await
        .expect("Failed to connect to database");
    Migrator::refresh(&db_connection).await.expect("Failed to run migrations");
    let config = listener_config(&rp).await;
    let consumer = shisha::messages::create_kafka_consumer(&config);
    let connection_clone = db_connection.clone();

    tokio::spawn(async move {
        shisha::messages::consume_events_task(consumer, &config, &connection_clone).await;
    });
    write_shisha_message(ShishaMessage {
        user: "TEST".into(),
        r#type: "upload".into(),
        ..ShishaMessage::default()
    }, &producer).await;
    write_shisha_message(ShishaMessage {
        user: "TEST".into(),
        r#type: "buy".into(),
        amount: Some(25),
        image_uuid: Some("123-345465353432-34234234".into()),
        ..ShishaMessage::default()
    }, &producer).await;
    let start = chrono::Utc::now();
    loop {
        let now = chrono::Utc::now();
        if (now - start).gt(&TimeDelta::seconds(10)) {
            panic!("Event not happened in required time")
        }
        let ev = QueryCore::last_10_events(&db_connection).await.expect("Failed to read messages from db");
        match ev.first() {
            Some(m) => {
                assert_eq!(m.actor, "TEST");
                assert_eq!(m.event_type, EventType::UPLOAD);
                break;
            }
            None => tokio::time::sleep(Duration::from_secs(1)).await
        }
    }
    let sum = QueryCore::count_money(&db_connection).await.expect("Failed to count money");
    assert_eq!(sum, 25)
}