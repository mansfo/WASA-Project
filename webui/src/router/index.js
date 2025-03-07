import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import ProfileView from '../views/ProfileView.vue'
import UserListView from '../views/UserListView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', redirect: '/session'},
		{path: '/users/:uid/stream', component: HomeView},
		{path: '/session', component: LoginView},
		{path: '/users/:uid/profile', component: ProfileView},
		{path: '/users/:uid/search/:searchedName', alias: ['/users/:uid/:searchedName', '/users/:uid/photos/:searchedName/likes'], component: UserListView},
	]
})

export default router