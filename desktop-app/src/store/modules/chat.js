
export default {
    state: () => ({
        openChatMessages:[],
        unreadMessages:{},
        openChat: null,
        chatUserList: [],
    }),
    getters: {

    },
    mutations: {
        updateOpenChatMessages(state, msgs) {
            console.log("Update open chat messages: ",  msgs)
            state.openChatMessages = msgs
        },
        openNewChat(state, openChat) {
            console.log("Open new chat: ", openChat)
            state.openChat = openChat
        },
        updateChatUserList(state, userList) {
            state.chatUserList = userList
        },
        updateUnreadMessagest(state, msgs) {
             console.log("Those are UNREAD",msgs)
            state.unreadMessages = msgs
        },
    },
    actions: {
        async fetchChatUserList({rootState, commit, dispatch}) {
            await dispatch("getMyUserID");
           
            const response = await fetch('http://localhost:8081/chatList?userId=' + rootState.id, {
                credentials: 'include'
            });
            const data = await response.json();
            commit("updateChatUserList", data.users);
        },
        async fecthChatMessages({state,commit}){
            const response = await fetch("http://localhost:8081/messages", {
                credentials: "include",
                method: "POST",
                body: JSON.stringify({
                    type: state.openChat.type,
                    receiverId: state.openChat.id
                })
                
            });
            const data = await response.json();
            let previousMessages = data.chatMessage ? data.chatMessage : [];
            console.log("MSg received from db: ", previousMessages)
            commit('updateOpenChatMessages', previousMessages)
        },
        async sendMessage({state,dispatch}, payload){
            const msgObj = {
                receiverId: state.openChat.id,
                content: payload,
                type: state.openChat.type
            };
             let response = await fetch("http://localhost:8081/newMessage", {
                body: JSON.stringify(msgObj),
                method: "POST",
                credentials: "include"
            });
            const data = await response.json();
            if (data.type == "Success"){
                dispatch("addNewChatMessage", data.chatMessage[0]);
                return true
            }else{
                return false
            }
        },
        addNewChatMessage({ commit, state }, payload) {
            let openMessages = state.openChatMessages;
            openMessages.push(payload)
            commit("updateOpenChatMessages", openMessages)
        },
        addUnreadMessage({ commit, state }, payload){
            let unreadMessages = state.unreadMessages
            if(unreadMessages[payload]){
               unreadMessages[payload] = unreadMessages[payload]+1
            }else{
                unreadMessages[payload] = 1
            }
            commit("updateUnreadMessagest", unreadMessages)
        },
         removeUnreadMessages({ commit, state }, payload){
            let unreadMessages = state.unreadMessages
            if(unreadMessages[payload]){
               delete unreadMessages[payload]
               commit("updateUnreadMessagest", unreadMessages)
            }
        }
    },
}