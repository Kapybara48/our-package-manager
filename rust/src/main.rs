use std::io::Write;

use serde::{Serialize, Deserialize};
use anyhow::{Context, Result};
use clap::{Parser, Subcommand};
use bytesize::ByteSize;

use crate::Commands::Install;

#[derive(Serialize, Deserialize, Debug)]
struct Release{
    id: i32,
    assets: Vec<Assets>,
}

#[derive(Serialize, Deserialize, Debug)]
struct Assets{
    name: String,
    digest: String,
    size: u64,
    browser_download_url: String,
}

#[derive(Parser)]
struct Cli{
    #[command(subcommand)]
    command: Commands,
}

#[derive(Subcommand)]
enum Commands{
    Install {url: String},
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>>{
    let cli = Cli::parse();
    match cli.command{
        Install{ url } => match install(&url).await{
            Ok(()) => println!("something went right"),
            Err(err) => eprintln!("everything went left in installation: {err:#}")
        }
    }
    Ok(())
}

async fn install(url: &str) -> Result<(), Box<dyn std::error::Error>>{
    let client = reqwest::Client::new();
    let response = client.get(url).header("User-Agent", "OUR").send().await.context("Happened when getting response")?;

    let release: Release = response.json().await.context("Happened when parsing to struct")?;
    
    for (index, asset) in release.assets.iter().enumerate()  {
        println!("{0}: {1} ({2})", index, asset.name, ByteSize(asset.size));
    }

    print!("Choose one to install: ");
    std::io::stdout().flush()?;
    let mut input = String::new();
    std::io::stdin().read_line(&mut input)?;

    let asset_index: usize = input.trim().parse()?;
    let selected_asset = release.assets.get(asset_index).ok_or("You chose wrong index")?;
    println!("{}", selected_asset.name);

    Ok(())
}
