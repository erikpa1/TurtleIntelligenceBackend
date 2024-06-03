use alloc::format;
use core::convert::Into;
use math::vec3::{Position};
use crate::inworld::InWorld;
use crate::stepper::Stepper;
use crate::tool_context::ToolsContext;

pub struct OutputEntity {
    pub name: String,
    pub position: Position,
}

impl OutputEntity {
    pub fn New() -> Self {
        OutputEntity {
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

impl InWorld for OutputEntity {
    fn GetName(&self) -> String {
        return self.name.clone();
    }

    fn GetType(&self) -> String {
        "input".into()
    }
}