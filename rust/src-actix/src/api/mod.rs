use actix_web::web::ServiceConfig;
use actix_web::{get, post, web, App, HttpResponse, HttpServer, Responder};
use std::sync::{Arc, Mutex};

use serde_json::{json, Value};

use super::app;


#[get("/api/simulation/state")]
async fn _GetAppState(data: web::Data<app::AppArcMutext>) -> impl Responder {
    let application_r = data.lock();

    if let Ok(applicaton) = application_r {
        let data = applicaton.to_json();

        return web::Json(data);
    }

    return web::Json(json!({"test": 1}));
}

pub fn mount(config: &mut ServiceConfig) {
    config.service(_GetAppState);
}
