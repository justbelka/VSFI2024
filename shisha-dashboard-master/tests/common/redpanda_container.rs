use testcontainers::{ContainerRequest, Image, ImageExt, TestcontainersError};
use testcontainers::core::{CmdWaitFor, ContainerPort, ContainerState, ExecCommand, WaitFor};

pub const REDPANDA_PORT: u16 = 9092;
pub const ADMIN_PORT: u16 = 9644;

#[derive(Debug)]
pub struct Redpanda {
    tag: String,
}

impl Redpanda {
    /// creates test container for specified tag
    pub fn for_tag(tag: String) -> ContainerRequest<Self> {
        ContainerRequest::from(Self { tag })
            .with_mapped_port(REDPANDA_PORT, ContainerPort::Tcp(REDPANDA_PORT))
            .with_mapped_port(ADMIN_PORT, ContainerPort::Tcp(ADMIN_PORT))
    }
}

#[allow(dead_code)]
impl Redpanda {
    /// A command to create new topic with specified number of partitions
    ///
    /// # Arguments
    ///
    /// - `topic_name` name of the topic to be created
    /// - `partitions` number fo partitions for given topic
    pub fn cmd_create_topic(topic_name: &str, partitions: i32) -> ExecCommand {
        log::debug!("cmd create topic [{}], with [{}] partition(s)", topic_name, partitions);
        // not the best ready_condition
        let container_ready_conditions = vec![
            WaitFor::StdErrMessage {
                message: "Create topics".into(),
            },
            WaitFor::Duration {
                length: std::time::Duration::from_secs(1),
            },
        ];

        ExecCommand::new(vec![
            String::from("rpk"),
            String::from("topic"),
            String::from("create"),
            String::from(topic_name),
            String::from("-p"),
            format!("{}", partitions),
        ]).with_cmd_ready_condition(CmdWaitFor::Duration {
            length: std::time::Duration::from_secs(1),
        })
            .with_container_ready_conditions(container_ready_conditions)
    }
}

impl Image for Redpanda {
    fn name(&self) -> &str {
        "redpandadata/redpanda"
    }

    fn tag(&self) -> &str {
        self.tag.as_str()
    }

    fn ready_conditions(&self) -> Vec<WaitFor> {
        vec![
            WaitFor::StdErrMessage {
                message: "Initialized cluster_id to ".into(),
            },
        ]
    }

    fn entrypoint(&self) -> Option<&str> {
        Some("sh")
    }

    fn cmd(&self) -> impl IntoIterator<Item=impl Into<std::borrow::Cow<'_, str>>> {
        vec![
            "-c",
            "/usr/bin/rpk redpanda start --mode dev-container --node-id 0 --set redpanda.auto_create_topics_enabled=false"
            ,
        ].into_iter()
    }

    fn expose_ports(&self) -> &[ContainerPort] {
        &[]
    }

    fn exec_after_start(&self, _: ContainerState) -> Result<Vec<ExecCommand>, TestcontainersError> {
        Ok(vec![])
    }
}
