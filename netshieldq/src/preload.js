const { contextBridge, ipcRenderer } = require('electron');

contextBridge.exposeInMainWorld('electron', {
    savePrivateKey: (privateKey, userId) => ipcRenderer.send('save-private-key', privateKey, userId),
    readPrivateKey: (userId) => ipcRenderer.sendSync('read-private-key', userId),
    checkFileExists: (fileName) => ipcRenderer.sendSync('check-file-exists', fileName),
    saveSharedSecret: (sharedSecretBase64, fileName) => ipcRenderer.send('save-shared-secret', sharedSecretBase64, fileName),
    readSharedSecret: (fileName) => ipcRenderer.sendSync('read-shared-secret', fileName),
    aesEncrypt: (data, key) => ipcRenderer.sendSync('aes-encrypt', data, key),
    aesDecrypt: (encryptedData, key) => ipcRenderer.sendSync('aes-decrypt', encryptedData, key),
    send: (channel, data) => {
        ipcRenderer.send(channel, data);
    },
    on: (channel, func) => {
        ipcRenderer.on(channel, (event, ...args) => func(...args));
    }
});
