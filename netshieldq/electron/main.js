const { app, BrowserWindow ,ipcMain} = require('electron');
const path = require('path');
process.env.NODE_ENV = 'production';
const fs = require('fs');
const crypto = require('crypto');

function createWindow() {
    const win = new BrowserWindow({
        width: 800,
        height: 600,
        webPreferences: {
            nodeIntegration: false,
            contextIsolation: true,
            preload: path.join(__dirname, '../src/preload.js'),
            icon: path.join(__dirname, 'icons8-3d-48.png'),

        },
        autoHideMenuBar: true,
        frame: false,
    });

    if (app.isPackaged) {
        win.loadFile('dist/index.html');
    } else {
        win.loadURL('http://localhost:8080/');
    }
    // win.webContents.openDevTools();

}
ipcMain.on('save-shared-secret', (event, sharedSecretBase64, fileName) => {
    const dirPath = path.join(app.getPath('appData'), 'netshieldq', 'keys');
    fs.mkdirSync(dirPath, { recursive: true }); // 创建目录，如果目录已存在则不会报错
    const filePath = path.join(dirPath, fileName);
    fs.writeFile(filePath, sharedSecretBase64, (err) => {
        if (err) {
            console.log('Failed to save the shared secret', err);
            event.returnValue = false;
        } else {
            console.log('Shared secret saved successfully!');
            event.returnValue = true;
        }
    });
});

ipcMain.on('read-shared-secret', (event, fileName) => {
    const dirPath = path.join(app.getPath('appData'), 'netshieldq', 'keys');
    const filePath = path.join(dirPath, fileName);
    fs.readFile(filePath, 'utf8', (err, data) => {
        if (err) {
            console.log('Failed to read the shared secret', err);
            event.returnValue = null;
        } else {
            console.log('Shared secret read successfully!');
            event.returnValue = data;
        }
    });
});
ipcMain.on('save-private-key', (event, privateKey, userId) => {
    const fileName = `privateKey_${userId}.txt`;
    const dirPath = path.join(app.getPath('appData'),'netshieldq', 'keys',`${userId}`);
    fs.mkdirSync(dirPath, { recursive: true });
    const filePath = path.join(dirPath, fileName);
    fs.writeFile(filePath, privateKey, (err) => {
        if (err) {
            console.log('Failed to save the private key', err);
        } else {
            console.log('Private key saved successfully!');
        }
    });
});
ipcMain.on('read-private-key', (event, userId) => {
    const fileName = `privateKey_${userId}.txt`;
    const dirPath = path.join(app.getPath('appData'), 'netshieldq', 'keys',`${userId}`);
    const filePath = path.join(dirPath, fileName);
    fs.readFile(filePath, 'utf8', (err, data) => {
        if (err) {
            console.log('Failed to read the private key', err);
            event.returnValue = null;
        } else {
            event.returnValue = data;
        }
    });
});

ipcMain.on('check-file-exists', (event, fileName) => {
    const dirPath = path.join(app.getPath('appData'), 'netshieldq', 'keys');
    const filePath = path.join(dirPath, fileName);
    fs.access(filePath, fs.constants.F_OK, (err) => {
        event.returnValue = !err;
    });
});
ipcMain.on('aes-encrypt', (event, data,key) => {
    const sharedSecretKey = Uint8Array.from(atob(key), c => c.charCodeAt(0));
    const cipher = crypto.createCipher('aes-256-cbc', sharedSecretKey);
    let crypted = cipher.update(data, 'utf8', 'hex');
    crypted += cipher.final('hex');
    event.returnValue = crypted;
});

ipcMain.on('aes-decrypt', (event, encryptedData, key) => {
    const sharedSecretKey = Uint8Array.from(atob(key), c => c.charCodeAt(0));

    const decipher = crypto.createDecipher('aes-256-cbc', sharedSecretKey);
    var encryptedData = encryptedData.toString();
    let decrypted = decipher.update(encryptedData, 'hex', 'utf8');
    decrypted += decipher.final('utf8');
    event.returnValue = decrypted;


});


ipcMain.on('minimize-window', () => {
    const win = BrowserWindow.getFocusedWindow();
    if (win) win.minimize();
});

ipcMain.on('maximize-window', () => {
    const win = BrowserWindow.getFocusedWindow();
    if (win) {
        if (win.isMaximized()) {
            win.unmaximize();
        } else {
            win.maximize();
        }
    }
});

ipcMain.on('close-window', () => {
    const win = BrowserWindow.getFocusedWindow();
    if (win) win.close();
});

app.on('ready', createWindow);

app.on('window-all-closed', () => {
    if (process.platform !== 'darwin') {
        app.quit();
    }
});

app.on('activate', () => {
    if (BrowserWindow.getAllWindows().length === 0) {
        createWindow();
    }
});
process.on('uncaughtException', (error) => {
    console.error('Uncaught Exception:', error);

    // 尝试局部重启窗口
    reloadWindow();
});

function reloadWindow() {
    const currentWindow = BrowserWindow.getFocusedWindow();
    if (currentWindow) {
        currentWindow.reload();
    }
}
