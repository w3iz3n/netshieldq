import Vue from 'vue';
import Vuex from 'vuex';

Vue.use(Vuex);

export default new Vuex.Store({
    state: {
        socket: null,
        historicalMessages: {},  // 存储历史消息
        messages: {},            // 存储实时消息
        currentUserID: null,     // 当前用户的 ID
        currentChat: null,       // 当前聊天好友的 ID
    },
    mutations: {
        SET_SOCKET(state, socket) {
            state.socket = socket;
        },
        CLEAR_MESSAGES(state) {
            state.messages = {};
            state.historicalMessages = {};
        },
        SET_CURRENT_USER_ID(state, userId) {
            state.currentUserID = userId;
        },
        SET_CURRENT_CHAT(state, friendId) {
            state.currentChat = friendId;
        },
        ADD_SEND_REALTIME_MESSAGE(state, message) {
            const friendId = state.currentChat;
            if (!state.realTimeMessages[friendId]) {
                state.realTimeMessages[friendId] = [];
            }
            state.realTimeMessages[friendId].push(message);
        },
        ADD_MESSAGE(state, message) {
            const key = message.receiverId === state.currentUserID ? message.senderId : message.receiverId;
            if (!state.messages[key]) {
                Vue.set(state.messages, key, []);
            }
            state.messages[key].push({
                id: message.id,
                senderId: message.senderId,
                receiverId: message.receiverId,
                text: message.text,
                timestamp: message.timestamp,
                isHistorical: false,
                messageType:message.messageType
            });
        },
        ADD_HISTORICAL_MESSAGE(state, message) {
            const key = message.ReceiverID === state.currentUserID ? message.SenderID : message.ReceiverID;
            if (!state.historicalMessages[key]) {
                Vue.set(state.historicalMessages, key, []);
            }
            state.historicalMessages[key].push({
                id: message.MessageId,
                senderId: message.SenderID,
                receiverId: message.ReceiverID,
                text: message.Content,
                messageType:message.MessageType,
                timestamp: message.Timestamp,
                isHistorical: message.IsHistorical
            });
        },
    },
    actions: {
        createSocket({ commit, dispatch }) {
            const jwt = localStorage.getItem('jwt');
            // console.log('createSocket')

            const socket = new WebSocket(`ws://localhost:3000/ws?token=${jwt}`);
            socket.onmessage = function(event) {
                console.log(event.data);
                if (event.data === "ping") {
                    console.log("Pong received");
                }else {
                const data = JSON.parse(event.data);
                dispatch('handleMessage', data);
                }
            };

            socket.onopen = function() {
                console.log("WebSocket connected.");
            };

            socket.onerror = function(error) {
                console.error("WebSocket error:", error);
            };

            socket.onclose = function() {
                console.log("WebSocket disconnected.");
                commit('SET_SOCKET', null);
            };

            commit('SET_SOCKET', socket);
        },
        closeSocket({ state, commit }) {
            if (state.socket) {
                state.socket.close();
                commit('SET_SOCKET', null);
            }
        },
        // 打开消息框
        initchat({ commit, dispatch }, friendId) {
            // console.log(friendId);
            commit('SET_CURRENT_USER_ID', localStorage.getItem('userid'));
            commit('SET_CURRENT_CHAT', friendId);
            dispatch('loadHistoricalMessages', friendId);
        },
        handleMessage({ commit, state }, message) {
            if (message.receiverId === parseInt(state.currentUserID) || message.senderId === parseInt(state.currentUserID)) {
                commit('ADD_MESSAGE', message);
            }
        },
        // 获取二者之间历史信息
        loadHistoricalMessages({ commit }, friendId ) {
            const userId = localStorage.getItem('userid');  // 从 localStorage 获取当前用户 ID
            if (!userId) {
                console.error('User ID is not set in localStorage.');
                return;
            }

            fetch(`/api/message/history/${userId}/${friendId}`)
                .then(response => {
                    if (!response.ok) {
                        throw new Error(`HTTP error! status: ${response.status}`);
                    }
                    return response.json();
                })
                .then(data => {
                    data.forEach(message => {
                        // console.log(message);
                        commit('ADD_HISTORICAL_MESSAGE', message);
                    });
                })
                .catch(error => console.error('Failed to load historical messages:', error));
        },
    }
});
