use clap::{Parser, Subcommand};

mod build;
mod cargo;
mod error;
mod git;
mod paths;

#[derive(Parser)]
#[command()]
struct Args {
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands {
    Install { url: String },
    Remove,
    Update,
}

fn main() {
    let args = Args::parse();

    match args.command {
        Commands::Install { url } => match install(&url) {
            Ok(()) => println!("successfully installed"),
            Err(error) => println!("{}", error),
        },
        Commands::Remove => println!("removing"),
        Commands::Update => println!("updating"),
    }
}

fn install(url: &str) -> Result<(), error::OurError> {
    println!("installing {}", url);

    let (status, package_dir) = git::clone(url)?;

    if !status.success() {
        println!("failed to clone");
        return Ok(());
    }

    println!("successfully cloned");

    let status = build::build(package_dir)?;

    if !status.success() {
        println!("failed to build");
        return Ok(());
    }

    println!("successfully built");

    Ok(())
}
