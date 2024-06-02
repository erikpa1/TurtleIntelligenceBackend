use alloc::format;
use core::convert::Into;
use math::vec3::{Position};
use crate::stepper::Stepper;
use crate::tool_context::ToolsContext;

pub struct Entity {
    pub name: String,
    pub position: Position,
}

impl Entity {
    pub fn New() -> Self {
        Entity {
            name: "".into(),
            position: Position::NewZero(),
        }
    }
    pub fn Step(&self, stepper: &Stepper, context: &mut ToolsContext) {
       // println!("Entity doing something random [{}]", context.expr.Execute(&"standard()".into()));
    }
    pub fn Init(&self) {
        println!(
            "File: {}, Line: {}, {}",
            file!(),
            line!(),
            format!("[{}] received init", &self.name)
        );
    }
}