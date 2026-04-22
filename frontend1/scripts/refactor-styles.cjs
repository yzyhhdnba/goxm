const fs = require('fs');
const path = require('path');

function walk(dir, callback) {
  const files = fs.readdirSync(dir);
  for (const f of files) {
    const filepath = path.join(dir, f);
    const stats = fs.statSync(filepath);
    if (stats.isDirectory()) {
      walk(filepath, callback);
    } else if (filepath.endsWith('.tsx') || filepath.endsWith('.ts')) {
      callback(filepath);
    }
  }
}

function processFile(filepath) {
  let content = fs.readFileSync(filepath, 'utf8');

  let oldContent = content;

  // Layout resets
  content = content.replace(/rounded-\[32px\]/g, 'rounded-md');
  content = content.replace(/rounded-3xl/g, 'rounded-lg');
  content = content.replace(/rounded-2xl/g, 'rounded-[6px]');
  
  // Backgrounds & effects
  content = content.replace(/bg-white\/80/g, 'bg-white');
  content = content.replace(/bg-paper\/80/g, 'bg-white');
  content = content.replace(/bg-paper\/90/g, 'bg-white');
  content = content.replace(/bg-[A-Za-z-]+\/[0-9]+/g, function(match) {
    // preserve explicit needed transparecies but generally strip AI misty stuff
    if (match.includes('ink') || match.includes('white') || match.includes('black')) return match;
    return match.split('/')[0]; 
  });
  content = content.replace(/bg-paper/g, 'bg-white');
  content = content.replace(/backdrop-blur-[a-z]+/g, '');
  content = content.replace(/shadow-float/g, 'shadow-sm border border-mist transition-all hover:shadow-md');
  
  // Fancy AI typography that Bilibili doesn't use
  content = content.replace(/uppercase tracking-\[0\.3em\]/g, 'font-medium tracking-wide');
  content = content.replace(/tracking-\[0\.3em\]/g, 'font-medium tracking-wide');
  content = content.replace(/tracking-\[0\.2em\]/g, 'font-medium');
  content = content.replace(/uppercase font-semibold tracking-wider/g, 'font-medium');

  if (oldContent !== content) {
    fs.writeFileSync(filepath, content);
    console.log('Updated:', filepath);
  }
}

walk(path.join(__dirname, '..', 'src'), processFile);
