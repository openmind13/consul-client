pub use anyhow::Result as AnyResult;
use serde::Deserialize;
use std::{
    any::Any,
    io::{BufReader, Error, Read},
};
use toml;

#[derive(Debug, Deserialize)]
pub struct Config {
    pub consul_address: String,
    pub token: String,
    pub http_listen_address: String,
}

impl Config {
    fn empty() -> Config {
        Config {
            consul_address: "".to_string(),
            token: "".to_string(),
            http_listen_address: "".to_string(),
        }
    }

    pub fn parse(path: &str) -> AnyResult<Config> {
        let mut file = std::fs::File::open(path)?;
        let mut buf = String::new();
        file.read_to_string(&mut buf)?;
        Ok(toml::from_str::<Config>(&buf)?)
    }
}
