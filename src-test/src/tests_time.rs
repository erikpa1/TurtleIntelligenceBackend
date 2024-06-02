use math::time::TimeExpresionExecutioner;

#[cfg(test)]
mod tests {
    use core::convert::Into;
    use super::*;

    #[test]
    fn it_works() {
        let result = TimeExpresionExecutioner::MakeFromMilis(0);
        assert_eq!(result, String::from("00:00"));

        let result = TimeExpresionExecutioner::MakeFromMilis(1000);
        assert_eq!(result, String::from("00:01"));

        let result = TimeExpresionExecutioner::MilisFromTimeString(&"00:01".into());
        assert_eq!(result, 1000);

        let result = TimeExpresionExecutioner::MilisFromTimeString(&"01:00".into());
        assert_eq!(result, 60 * 1000);
    }
}
