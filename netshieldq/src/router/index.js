import Vue from 'vue'
import Router from 'vue-router'
import HomePage from '@/components/HomePage.vue';
import IntraScan from "@/components/Intra/IntraLogin.vue";
import LogIn from "@/components/public/LogIn.vue";
import ChatApp from '@/views/ChatApp.vue';
import friendRequest from "@/components/friendRequest.vue";
import FriendShip from "@/views/FriendShip.vue";
import FileManage from "@/components/FileManage.vue";
import KyberEnc from "@/components/KyberEnc.vue";
import Register from "@/components/public/RegisterFunc.vue";
Vue.use(Router)

export default new Router({
    mode:'history',
    routes: [
        {
            path: '/home',
            name: 'HomePage',
            component: HomePage
        },
        {
            path:'/intrascan',
            name:'intraScan',
            component:IntraScan
        },
        {
            path:'/plogin',
            name:"LogIn",
            component:LogIn
        },
        {
            path: '/chat',
            name: 'chat',
            component: ChatApp,
            // meta: {
            //     requiresAuth: true // 假设聊天页面需要用户登录认证
            // }
        },
        {
            path:'/sendrequest',
            name:'sendRequest',
            component:friendRequest
        },{
            path:"/friendship",
            name:"friendShip",
            component : FriendShip
        },{
            path:"/files",
            name:"file",
            component:FileManage
         },{
            path:'/kyber',
            name:"KyberEnc",
            component:KyberEnc
        },{
            path:"/register",
            name:"RegisterFunc",
            component:Register
        }
    ]
})
const VueRouterPush = Router.prototype.push
Router.prototype.push = function push (to) {
    return VueRouterPush.call(this, to).catch(err => err)
}
