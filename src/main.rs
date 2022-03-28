pub use anyhow::Result as AnyResult;
use consul;
use consul::agent::Agent;
use consul::agent::AgentCheck;
use consul::catalog::Catalog;
use consul::catalog::CatalogRegistration;
use consul::session::Session;
use consul::session::SessionEntry;
use consul::WriteOptions;
use std::any::Any;
use std::collections::HashMap;
use std::fmt::Write;

mod config;

// fn main() {
//     println!("Consul-rs");

//     let cfg_path = {
//         match std::env::var("CFG_PATH") {
//             Ok(cfg_path) => cfg_path,
//             Err(_) => {
//                 println!("Specify env CFG_PATH");
//                 std::process::exit(0);
//             }
//         }
//     };

//     let config = match config::Config::parse(&cfg_path) {
//         Ok(config) => config,
//         Err(e) => {
//             println!("Failed to parse config {}", e);
//             std::process::exit(-1);
//         }
//     };

//     let mut consul_config = consul::Config::new().unwrap();
//     consul_config.address = config.consul_address;
//     consul_config.token = Some(config.token);

//     let mut client = consul::Client::new(consul_config);

//     // let agent_check = consul::agent::AgentCheck::default();

//     // let members = client.members(false).unwrap();
//     // println!("{:?}", members);

//     let (services, other) = client.services(None).unwrap();
//     println!("{:?}", services);

//     println!("{:?}", other);
// }

use consulrs::api::check::common::AgentServiceCheckBuilder;
use consulrs::api::service::requests::RegisterServiceRequest;
use consulrs::client::{ConsulClient, ConsulClientSettingsBuilder};
use consulrs::service;
use std::convert::TryInto;
use tokio::runtime;

async fn run() -> AnyResult<()> {
    let cfg_path = {
        match std::env::var("CFG_PATH") {
            Ok(cfg_path) => cfg_path,
            Err(_) => {
                println!("Specify env CFG_PATH");
                std::process::exit(0);
            }
        }
    };

    let config = match config::Config::parse(&cfg_path) {
        Ok(config) => config,
        Err(e) => {
            println!("Failed to parse config {}", e);
            std::process::exit(-1);
        }
    };

    let client = ConsulClient::new(
        ConsulClientSettingsBuilder::default()
            .address(&config.consul_address)
            .build()
            .unwrap(),
    )
    .unwrap();

    let check = AgentServiceCheckBuilder::default()
        .name("rust-test-check")
        .interval("5s")
        .build()
        .unwrap();

    // let register_request = RegisterServiceRequest::builder()
    //     .address(config.consul_address)
    //     // .port(1234)
    //     .check(check);

    service::register(
        &client,
        "rust-test-service-dev",
        Some(
            RegisterServiceRequest::builder()
                .address(config.consul_address)
                .port(1234)
                .check(
                    AgentServiceCheckBuilder::default()
                        .name("rust-test-http-check")
                        .interval("10s")
                        .http("http://myservice.lab.com/health")
                        .status("passing")
                        .build()
                        .unwrap(),
                ),
        ),
    )
    .await;

    Ok(())
}

fn main() {
    println!("Consul-rs");

    let rt = runtime::Builder::new_multi_thread()
        .worker_threads(2)
        .enable_all()
        .build()
        .unwrap();
    rt.block_on(run()).unwrap();
}
