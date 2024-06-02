use math::expr::MathExpresionExecutioner;

pub struct ToolsContext {
    pub expr: MathExpresionExecutioner,
}

impl ToolsContext {
    pub fn New() -> Self {
        ToolsContext {
            expr: MathExpresionExecutioner::New()
        }
    }
}
