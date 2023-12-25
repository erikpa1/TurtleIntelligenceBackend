use std::collections::HashMap;

use serde_json::{json, Value};

pub struct Warehouse {
    resmap: HashMap<String, i32>,
}

impl Warehouse {
    pub fn new() -> Warehouse {
        Warehouse {
            resmap: HashMap::new(),
        }
    }

    pub fn set_resource(&mut self, name: &str, initial_value: i32) {
        self.resmap.insert(name.into(), initial_value);
    }

    pub fn add_resouce(&mut self, name: &str, inc_value: i32) {}

    pub fn lower_resource(&mut self, name: &str, dec_value: i32) {
        let value = self.resmap.get(name).copied().unwrap_or(0);
        let mut new_value = value - dec_value;

        if new_value < 0 {
            new_value = 0
        }

        self.resmap.insert(name.into(), new_value);
    }

    pub fn get_resource_count(&self, name: &str) -> i32 {
        self.resmap.get(name).copied().unwrap_or(0)
    }

    pub fn to_json(&self) -> Value {
        json!({
            "resources": self.resmap
        })
    }
}
