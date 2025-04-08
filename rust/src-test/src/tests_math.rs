use math::expr::MathExpresionExecutioner;

#[cfg(test)]
mod tests {
    use core::convert::Into;
    use super::*;

    #[test]
    fn it_works() {
        let mut expr = MathExpresionExecutioner::New();

        let result = expr.Execute(&"standard()".into());
        assert_eq!(result < 1.0, true);

        let result = expr.Execute(&"uniform(1, 10)".into());
        assert_eq!(result > 0.0 && result <= 10.0, true);
    }
}
