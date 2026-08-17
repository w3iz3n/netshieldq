const fs = require('fs');
const yaml = require('js-yaml');
const path = require('path');

function read_config(){
    const filePath = path.join(__dirname, '../../config/config.yaml');
    const yamlData = fs.readFileSync(filePath, 'utf8');
    const config = yaml.load(yamlData);
    return config;
}
module.exports = read_config;
