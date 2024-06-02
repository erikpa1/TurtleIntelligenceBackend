use alloc::string::String;


pub struct TimeExpresionExecutioner {}

impl TimeExpresionExecutioner {
    pub fn New() -> Self {
        Self {

        }
    }
    pub fn MakeFromMilis(millis: u64) -> String {
        if millis == 0 {
            return String::from("00:00");
        }

        let seconds = millis / 1000;
        // Calculate days, hours, minutes, and remaining seconds
        let (mut days, mut seconds) = (seconds / 86400, seconds % 86400);
        let (mut hours, mut seconds) = (seconds / 3600, seconds % 3600);
        let (mut minutes, seconds) = (seconds / 60, seconds % 60);

        // Initialize the time components as strings
        let days_str = if days == 0 { String::new() } else { format!("{:02}:", days) };
        let hours_str = if hours == 0 { String::new() } else { format!("{:02}:", hours) };
        let minutes_str = if minutes == 0 { String::new() } else { format!("{:02}:", minutes) };
        let seconds_str = format!("{:02}", seconds);

        // Construct the time string based on the duration
        let time_string = if days > 0 {
            format!("{}{}{}{}", days_str, hours_str, minutes_str, seconds_str)
        } else if hours > 0 {
            format!("{}{}{}", hours_str, minutes_str, seconds_str)
        } else if minutes > 0 {
            format!("{}{}", minutes_str, seconds_str)
        } else {
            format!("00:{}", seconds_str)
        };

        time_string
    }

    pub fn SecondsFromTimeString(time_string: &String) -> i64 {
        return Self::MilisFromTimeString(time_string) / 1000
    }
    pub fn MilisFromTimeString(time_string: &String) -> i64 {
        // Split the time string into its components
        let mut components = time_string.split(":").collect::<Vec<&str>>();

        // Initialize variables for days, hours, minutes, and seconds
        let mut days = 0;
        let mut hours = 0;
        let mut minutes = 0;
        let mut seconds = 0;

        // Pop values from the end of the iterator until it's empty
        if let Some(sec) = components.pop() {
            seconds = sec.parse::<i64>().unwrap_or(0);
        }
        if let Some(min) = components.pop() {
            minutes = min.parse::<i64>().unwrap_or(0);
        }
        if let Some(hour) = components.pop() {
            hours = hour.parse::<i64>().unwrap_or(0);
        }
        if let Some(day) = components.pop() {
            days = day.parse::<i64>().unwrap_or(0);
        }

        // Calculate the total milliseconds
        let total_millis = (days * 24 * 60 * 60 + hours * 60 * 60 + minutes * 60 + seconds) * 1000;

        total_millis
    }
}