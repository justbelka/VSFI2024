use std::borrow::Cow;
use testcontainers::core::{ContainerPort, WaitFor};
use testcontainers::{ContainerAsync, Image};

const MYSQL_DATABASE: &str = "test_database";
const MYSQL_USER: &str = "test_user";
const MYSQL_PASSWORD: &str = "test_password_123";

const MYSQL_PORT: ContainerPort = ContainerPort::Tcp(3306);

pub struct Mysql {
    pub tag: String,
}

impl Mysql {
    pub fn for_tag(tag: String) -> Mysql {
        Mysql { tag }
    }
    pub async fn url(container: &ContainerAsync<Mysql>) -> String {
        format!("mysql://{}:{}@localhost:{}/{}",
                MYSQL_USER,
                MYSQL_PASSWORD,
                container.get_host_port_ipv4(MYSQL_PORT).await.expect("Failed to get mysql port from container"),
                MYSQL_DATABASE)
    }
}

impl Image for Mysql {
    fn name(&self) -> &str {
        "mysql"
    }

    fn tag(&self) -> &str {
        self.tag.as_str()
    }

    fn ready_conditions(&self) -> Vec<WaitFor> {
        vec![
            WaitFor::message_on_stderr("X Plugin ready for connections. Bind-address"),
            WaitFor::message_on_stderr("/usr/sbin/mysqld: ready for connections."),
        ]
    }

    fn env_vars(&self) -> impl IntoIterator<Item=(impl Into<Cow<'_, str>>, impl Into<Cow<'_, str>>)> {
        [
            ("MYSQL_DATABASE", MYSQL_DATABASE),
            ("MYSQL_USER", MYSQL_USER),
            ("MYSQL_PASSWORD", MYSQL_PASSWORD),
            ("MYSQL_ROOT_PASSWORD", MYSQL_PASSWORD),
        ]
    }

    fn expose_ports(&self) -> &[ContainerPort] {
        &[MYSQL_PORT]
    }
}