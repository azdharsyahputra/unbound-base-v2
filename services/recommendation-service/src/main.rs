use axum::{routing::get, Json, Router};
use serde_json::json;
use std::{env, net::SocketAddr};
use tokio::net::TcpListener;
use dotenvy::dotenv;

#[tokio::main]
async fn main() {
    dotenv().ok();

    let service_name = env::var("SERVICE_NAME").unwrap_or_else(|_| "unknown-service".to_string());
    let port = env::var("PORT").unwrap_or_else(|_| "8080".to_string());
    let addr: SocketAddr = format!("0.0.0.0:{port}").parse().unwrap();

    let listener = TcpListener::bind(addr).await.unwrap();
    let app = Router::new().route("/", get(root));

    println!("🚀 {service_name} running on 0.0.0.0:{port}");
    axum::serve(listener, app).await.unwrap();
}

async fn root() -> Json<serde_json::Value> {
    let service_name = std::env::var("SERVICE_NAME").unwrap_or_else(|_| "unknown-service".to_string());
    Json(json!({
        "service": service_name,
        "status": "running ✅"
    }))
}
