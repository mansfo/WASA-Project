import {createApp, reactive} from 'vue'
import App from './App.vue'
import router from './router'
import axios from './services/axios.js';
import ErrorMsg from './components/ErrorMsg.vue'
import LoadingSpinner from './components/LoadingSpinner.vue'
import Navbar from './components/Navbar.vue'
import Comment from './components/Comment.vue'
import CommentModal from './components/CommentModal.vue'
import Photo from './components/Photo.vue'
import DropDownMenu from './components/DropDownMenu.vue'
import User from './components/User.vue'
import LoginView from './views/LoginView.vue'
import HomeView from './views/HomeView.vue'
import ProfileView from './views/ProfileView.vue'
import UserListView from './views/UserListView.vue'

import './assets/dashboard.css'
import './assets/main.css'

const app = createApp(App)
app.config.globalProperties.$axios = axios;
app.component("ErrorMsg", ErrorMsg);
app.component("LoadingSpinner", LoadingSpinner);
app.component("Navbar", Navbar);
app.component("Comment", Comment);
app.component("CommentModal", CommentModal);
app.component("Photo", Photo);
app.component("DropDownMenu", DropDownMenu);
app.component("User", User)
app.use(router)
app.mount('#app')