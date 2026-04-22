import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const dirPath = "/Users/hhd/Desktop/test/goxm/frontend1/src";

function walk(dir) {
  let results = [];
  const list = fs.readdirSync(dir);
  list.forEach(file => {
    file = dir + '/' + file;
    const stat = fs.statSync(file);
    if (stat && stat.isDirectory()) { 
      results = results.concat(walk(file));
    } else {
      if (file.endsWith('.tsx') || file.endsWith('.ts')) {
        results.push(file);
      }
    }
  });
  return results;
}

const replacements = [
  [/bg-white\/80/g, 'bg-white'],
  [/bg-white\/70/g, 'bg-white'],
  [/bg-white\/85/g, 'bg-white'],
  [/bg-white\/75/g, 'bg-white'],
  [/bg-paper\/80/g, 'bg-paper'],
  [/rounded-\[32px\]/g, 'rounded-md'],
  [/rounded-\[28px\]/g, 'rounded-md'],
  [/rounded-\[24px\]/g, 'rounded-md'],
  [/rounded-3xl/g, 'rounded-md'],
  [/rounded-2xl/g, 'rounded-md'],
  [/shadow-float/g, 'shadow-sm'],
  [/shadow-2xl/g, 'shadow-md'],
  [/shadow-xl/g, 'shadow-md'],
  [/shadow-lg/g, 'shadow-sm'],
];

const files = walk(dirPath);
files.forEach(file => {
  let content = fs.readFileSync(file, 'utf8');
  let newContent = content;
  replacements.forEach(([regex, repl]) => {
    newContent = newContent.replace(regex, repl);
  });
  if (content !== newContent) {
    fs.writeFileSync(file, newContent, 'utf8');
    console.log('Updated', file);
  }
});

console.log('Done');
