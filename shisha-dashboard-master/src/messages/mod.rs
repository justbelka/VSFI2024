use metrics::counter;
use rdkafka::{ClientConfig, Message};
use rdkafka::config::RDKafkaLogLevel;
use rdkafka::consumer::{Consumer, StreamConsumer};
use sea_orm::{DatabaseConnection};
use serde::{Deserialize, Serialize};
use crate::service::query::Query;

pub struct ListenerConfig {
    pub boostrap_servers: String,
    pub topic: String,
}

#[derive(Serialize, Deserialize, Debug, Default)]
pub struct ShishaMessage {
    pub user: String,
    pub r#type: String,
    pub image_uuid: Option<String>,
    pub target: Option<String>,
    pub amount: Option<i64>,
}

pub fn create_kafka_consumer(config: &ListenerConfig) -> StreamConsumer {
    ClientConfig::new()
        .set("group.id", "shisha-dashboard-group")
        .set("bootstrap.servers", config.boostrap_servers.as_str())
        .set("enable.partition.eof", "false")
        .set("allow.auto.create.topics", "false")
        .set("session.timeout.ms", "6000")
        .set("enable.auto.commit", "true")
        .set("enable.auto.offset.store", "false")
        .set_log_level(RDKafkaLogLevel::Debug)
        .create()
        .expect("Failed to create kafka consumer")
}

pub async fn consume_events_task(con: StreamConsumer, config: &ListenerConfig, db: &DatabaseConnection) {
    con.subscribe(&[config.topic.as_str()]).expect("Failed to subscribe!");
    loop {
        match con.recv().await {
            Err(e) => log::error!("Kafka error: {}", e),
            Ok(m) => {
                let Some(payload) = m.payload() else {
                    log::error!("Failed to read message payload");
                    continue;
                };
                let message: ShishaMessage = match serde_json::from_slice(payload) {
                    Ok(res) => res,
                    Err(e) => {
                        log::error!("Failed to deserialize shisha event: {e}");
                        continue;
                    }
                };
                log::info!("Got message from shisha: {message:?}");
                match Query::insert_message(message, db).await {
                    Ok(_) => {
                        counter!("saved_events_count").increment(1);
                        let _ = con.store_offset_from_message(&m).inspect_err(|e| log::error!("Failed to commit offset: {e}"));
                    }
                    Err(e) => {
                        log::error!("Failed to write message in db: {e}, not commiting offset")
                    }
                }
            }
        }
    }
}

