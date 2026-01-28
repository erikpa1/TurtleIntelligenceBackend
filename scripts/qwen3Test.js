#!/usr/bin/env bun

/**
 * Ollama Qwen3-VL Image Analysis Script
 * Sends an image to Ollama's Qwen3-VL model for analysis
 */

const IMAGE_PATH = "/Users/erik.palencik/Desktop/testDeepsek/warehouse.png";
const OLLAMA_API_URL = "http://localhost:11434/api/generate";
const MODEL_NAME = "qwen3-vl:8b";

async function analyzeImage() {
  try {
    console.log("🖼️  Reading image from:", IMAGE_PATH);
    
    // Read the image file
    const imageFile = Bun.file(IMAGE_PATH);
    const imageBuffer = await imageFile.arrayBuffer();
    const base64Image = Buffer.from(imageBuffer).toString('base64');
    
    console.log("📤 Sending request to Ollama...\n");
    
    // Prepare the request payload
    const payload = {
      model: MODEL_NAME,
      prompt: "Describe this image in detail. What do you see?",
      images: [base64Image],
      stream: false
    };
    
    // Send POST request to Ollama (matching the cURL example)
    const response = await fetch(OLLAMA_API_URL, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify(payload)
    });
    
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Ollama API error: ${response.status} ${response.statusText}\n${errorText}`);
    }
    
    const result = await response.json();
    
    // Display the result
    console.log("🤖 Ollama Response:");
    console.log("═".repeat(80));
    console.log(result.response);
    console.log("═".repeat(80));
    console.log("\n✅ Analysis complete!");
    
    // Display metadata if available
    if (result.total_duration) {
      console.log("\n📊 Metadata:");
      console.log(`- Model: ${result.model || MODEL_NAME}`);
      console.log(`- Total duration: ${(result.total_duration / 1e9).toFixed(2)}s`);
      if (result.load_duration) {
        console.log(`- Load duration: ${(result.load_duration / 1e9).toFixed(2)}s`);
      }
    }
    
  } catch (error) {
    console.error("❌ Error:", error.message);
    
    if (error.code === "ENOENT") {
      console.error("\n💡 Tip: Make sure the image file exists at the specified path.");
    } else if (error.message.includes("ECONNREFUSED") || error.message.includes("fetch failed")) {
      console.error("\n💡 Tip: Make sure Ollama is running. Start it with: ollama serve");
    } else if (error.message.includes("model") || error.message.includes("not found")) {
      console.error("\n💡 Tip: Make sure the model is installed. Run: ollama pull qwen3-vl");
    }
    
    process.exit(1);
  }
}

// Run the script
console.log("🚀 Starting Ollama Image Analysis\n");
analyzeImage();