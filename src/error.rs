pub enum OurError {
    MissingHomeEnv,
    Io(std::io::Error),
}

impl From<std::io::Error> for OurError {
    fn from(error: std::io::Error) -> Self {
        OurError::Io(error)
    }
}
