use std::time::Duration;
use rdkafka::ClientConfig;
use rdkafka::producer::{FutureProducer, FutureRecord};
use testcontainers::ContainerAsync;
use testcontainers::runners::AsyncRunner;
use shisha::messages::{ListenerConfig, ShishaMessage};
use crate::common::mysql_container::Mysql;
use crate::common::redpanda_container::{Redpanda, REDPANDA_PORT};

pub mod redpanda_container;
pub mod mysql_container;

pub const TOPIC: &str = "shisha-events-topic";


pub async fn setup() -> (ContainerAsync<Redpanda>, ContainerAsync<Mysql>, FutureProducer) {
    env_logger::builder().is_test(true).init();
    let redpanda = setup_redpanda().await;
    let mysql = Mysql::for_tag("8.4.0".into()).start().await.expect("Failed to start MySQL");
    let producer = producer_setup(&redpanda).await;
    (redpanda, mysql, producer)
}

pub async fn db_url(mysql: &ContainerAsync<Mysql>) -> String {
    Mysql::url(&mysql).await
}

pub async fn listener_config(rp: &ContainerAsync<Redpanda>) -> ListenerConfig {
    ListenerConfig {
        boostrap_servers: bootstrap_servers(rp).await,
        topic: TOPIC.into(),
    }
}

pub async fn write_shisha_message(message: ShishaMessage, producer: &FutureProducer) {
    let json_message = serde_json::to_string(&message).expect("Failed to serialize message");
    let record = FutureRecord::to(TOPIC).key("TEST").payload(&json_message);
    producer.send(record, Duration::from_secs(0)).await.expect("Failed to send message to red panda");
}

async fn bootstrap_servers(rp: &ContainerAsync<Redpanda>) -> String {
    format!("localhost:{}",
            rp.get_host_port_ipv4(REDPANDA_PORT).await.expect("Failed to get red panda port"))
}

async fn producer_setup(rp: &ContainerAsync<Redpanda>) -> FutureProducer {
    ClientConfig::new()
        .set("bootstrap.servers", bootstrap_servers(rp).await)
        .set("message.timeout.ms", "5000")
        .create()
        .expect("Failed to create producer")
}

async fn setup_redpanda() -> ContainerAsync<Redpanda> {
    let rp = Redpanda::for_tag("v23.3.18".to_string()).start().await.expect("Failed to start Red Panda");
    rp.exec(Redpanda::cmd_create_topic(TOPIC, 1)).await.expect("Failed to create topic");
    rp
}