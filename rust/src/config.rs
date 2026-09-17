use serde::Deserialize;

#[derive(Deserialize)]
struct Config {
    package: Package,
    version: Version,
    source: Source,
    git: Git,
    build: Build,
    install: Install,
}

#[derive(Deserialize)]
struct Package {
    name: String,
    path: String,
}

#[derive(Deserialize)]
struct Version {
    commit: String,
}

#[derive(Deserialize)]
struct Source {
    url: String,
}

#[derive(Deserialize)]
struct Git {
    clone_depth: u32,
    branch: String,
}

#[derive(Deserialize)]
struct Build {
    command: String,
    args: Vec<String>,
}

#[derive(Deserialize)]
struct Install {
    binary_source: String,
    binary_destination: String,
}
