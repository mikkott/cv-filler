const puppeteer = require('puppeteer');
const path = require('path');

async function convertHtmlToPdf(htmlFilePath, pdfFilePath) {
    const browser = await puppeteer.launch();
    const page = await browser.newPage();
    
    const absoluteHtmlFilePath = path.resolve(htmlFilePath);
    await page.goto(`file://${absoluteHtmlFilePath}`, { waitUntil: 'networkidle0' });
    
    await page.pdf({ path: pdfFilePath, format: 'A4' });
    
    await browser.close();
}

async function convertHtmlToPng(htmlFilePath, pngFilePath) {
    const browser = await puppeteer.launch();
    const page = await browser.newPage();
    
    const absoluteHtmlFilePath = path.resolve(htmlFilePath);
    await page.goto(`file://${absoluteHtmlFilePath}`, { waitUntil: 'networkidle0' });
    
    await page.screenshot({ path: pngFilePath });
    
    await browser.close();
}

const htmlFilePath = 'output.html';
const pdfFilePath = 'output.pdf';
const pngFilePath = 'output.png';

convertHtmlToPdf(htmlFilePath, pdfFilePath)
    .then(() => {
        console.log('PDF conversion completed successfully!');
    })
    .catch((error) => {
        console.error('PDF conversion failed:', error);
    });

convertHtmlToPng(htmlFilePath, pngFilePath)
    .then(() => {
        console.log('PNG conversion completed successfully!');
    })
    .catch((error) => {
        console.error('PNG conversion failed:', error);
    });