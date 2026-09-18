use std::fmt::Formatter;

pub enum OurError {
    MissingHomeEnv,
    CloneFailed,
    BuildFailed,
    UnknownProjectType,
    Io(std::io::Error),
    TomlParse(toml::de::Error),
}

impl From<std::io::Error> for OurError {
    fn from(error: std::io::Error) -> Self {
        OurError::Io(error)
    }
}

impl From<toml::de::Error> for OurError {
    fn from(error: toml::de::Error) -> Self {
        OurError::TomlParse(error)
    }
}

impl std::fmt::Display for OurError {
    fn fmt(&self, f: &mut Formatter) -> std::fmt::Result {
        match self {
            OurError::MissingHomeEnv => write!(f, "Could not find your home directory"),
            OurError::CloneFailed => write!(f, "Failed to clone"),
            OurError::BuildFailed => write!(f, "Failed to build"),
            OurError::UnknownProjectType => write!(f, "Project type is unknown"),
            OurError::Io(error) => write!(f, "I/O error: {}", error),
            OurError::TomlParse(error) => write!(f, "Error while parsing toml: {}", error),
        }
    }
}
