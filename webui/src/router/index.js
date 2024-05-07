import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import LoginView from '../views/LoginView.vue'
import ProfileView from '../views/ProfileView.vue'
import MainstreamView from '../views/StreamView.vue'

const router = createRouter({
	history: createWebHashHistory(import.meta.env.BASE_URL),
	routes: [
		{path: '/', component: LoginView},
		{path: '/:username/profile', component: ProfileView},
		{path: '/:username/photostream', component: MainstreamView},
		{path: '/some/:id/link', component: HomeView},

	]
})

export default router
