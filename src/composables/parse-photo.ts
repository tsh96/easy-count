import { GoogleGenerativeAI } from "@google/generative-ai";
import { settings } from "./settings";

export async function parsePhoto(dataUrl: string) {
    const genAI = new GoogleGenerativeAI(settings.value.geminiApiKey);
  // The Gemini 1.5 models are versatile and work with both text-only and multimodal prompts

  const model = genAI.getGenerativeModel({
    model: settings.value.geminiModel,
    systemInstruction: `You are a specialized data extraction assistant. Your task is to extract data from photos containing tabular information. The photo will have the following columns: Date, Invoice No, Name, Amount (split into two columns), Date, Cheque, and Amount (split into two columns). Combine the two 'Amount' columns using a dot (\`.\`) to form a single decimal number (e.g., \`123\` and \`45\` become \`123.45\`). Format the extracted data as tab-separated values (TSV), ensuring each value is separated by a tab (\`\t\`). Exclude the header row. For example:
\`\`\`
2024-10-14	011452	Jabatan Pertanian Bhg SA	2257.70	2024-10-25	EPT	
2024-10-15	011453	Hospital SA	3000.00			
	011454	Thai Chiang	466.80			
\`\`\``,
  });

  const generationConfig = {
    temperature: 0,
    topP: 0.95,
    topK: 40,
    maxOutputTokens: 8192,
    responseMimeType: "text/plain",
  };

  const chatSession = model.startChat({
    generationConfig,
    history: [],
  })

  const generatedContent = await chatSession.sendMessage([{ inlineData: { data: dataUrl, mimeType: 'image/webp' } }, { text: '' }]);

  const response = generatedContent.response;
  // remove ``` from the start and end of the response if present
  const result = response.text().replace(/^```.*\n|```$/g, '')

  const usageMetadata = response.usageMetadata
  const cost = (usageMetadata?.candidatesTokenCount || 0) * 1.05 * 1e-6
    + (usageMetadata?.promptTokenCount || 0) * 0.35 * 1e-6

  return { cost, result }
}