use actix_files as fs;
use actix_web::middleware::Logger;
use actix_web::web::ServiceConfig;
use actix_web::{get, post, web, App, HttpResponse, HttpServer, Responder};

//Toto musi byt kvoli start v App
use actix::prelude::*;

mod app;
mod warehouse;

#[get("/api/test")]
async fn _test_function() -> impl Responder {
    HttpResponse::Ok().body("ok")
}

pub fn mount(reference: &mut ServiceConfig) {
    reference.service(_test_function);
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let mut app = app::App::new();
    app.start();

    HttpServer::new(move || {
        App::new()
            .wrap(Logger::default())
            .app_data(web::FormConfig::default())
            .configure(mount)
            .service(
                fs::Files::new("/", "static")
                    .index_file("index.html")
                    .show_files_listing(),
            )
    })
    .bind(("0.0.0.0", 5000))?
    .run()
    .await
}
