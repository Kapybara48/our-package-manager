use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Debug)]
struct Release{
    id: i32,
    assets: Vec<Assets>,
}

#[derive(Serialize, Deserialize, Debug)]
struct Assets{
    name: String,
    digest: String,
    browser_download_url: String,
}

const URL:&str = "";

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>>{
    let client = reqwest::Client::new();
    let response = client.get(URL).header("User-Agent", "OUR").send().await?;

    let release: Release = response.json().await?;
    let pretty = serde_json::to_string_pretty(&release)?;
    println!("{}", pretty);

    Ok(())
}
