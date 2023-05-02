import { useState } from 'react';
import { Redirect, Route, Switch } from 'react-router-dom';
import routes from 'src/horizon/routes';

// Chakra imports
import { Box, useColorModeValue } from '@chakra-ui/react';

// Layout components
import { SidebarContext } from 'src/horizon/contexts/SidebarContext';

// Custom Chakra theme
export default function Auth() {
	// states and functions
	const [ toggleSidebar, setToggleSidebar ] = useState(false);
	const getRoute = () => {
		return window.location.pathname !== '/auth/full-screen-maps';
	};
	// const getRoutesold = (
	// 	routes: RoutesType[]
	// ): any => {
	// 	return routes.map((route: RoutesType, key: any) => {
	// 		// optional props
	// 		if (route.layout === '/auth') {
	// 			return <Route path={route.layout + route.path} component={route.component} key={key} />;
	// 		} else {
	// 			return null;
	// 		}
	// 	});
	// };
	// const getActiveRoute = (
	// 	routes: RoutesType[]
	// ): any => {
	// 	let activeRoute = 'Default Brand Text';
	// 	for (let i = 0; i < routes.length; i++) {
	// 		if (routes[i].collapse && routes[i].items) {
	// 			let collapseActiveRoute = getActiveRoute(routes[i].items);
	// 			if (collapseActiveRoute !== activeRoute) {
	// 				return collapseActiveRoute;
	// 			}
	// 		} else {
	// 			if (window.location.href.indexOf(routes[i].layout + routes[i].path) !== -1) {
	// 				return routes[i].name;
	// 			}
	// 		}
	// 	}
	// 	return activeRoute;
	// };
	// const getActiveNavbar = (
	// 	routes: RoutesType[]
	// ): any => {
	// 	let activeNavbar = false;
	// 	for (let i = 0; i < routes.length; i++) {
	// 		if (routes[i].collapse) {
	// 			let collapseActiveNavbar = getActiveNavbar(routes[i].items);
	// 			if (collapseActiveNavbar !== activeNavbar) {
	// 				return collapseActiveNavbar;
	// 			}
	// 		} else {
	// 			if (window.location.href.indexOf(routes[i].layout + routes[i].path) !== -1) {
	// 				return routes[i].secondary;
	// 			}
	// 		}
	// 	}
	// 	return activeNavbar;
	// };
	// const getActiveNavbarText = (
	// 	routes: RoutesType[]
	// ): any => {
	// 	let activeNavbar = false;
	// 	for (let i = 0; i < routes.length; i++) {
	// 		if (routes[i].collapse) {
	// 			let collapseActiveNavbar = getActiveNavbarText(routes[i].items);
	// 			if (collapseActiveNavbar !== activeNavbar) {
	// 				return collapseActiveNavbar;
	// 			}
	// 		} else {
	// 			if (window.location.href.indexOf(routes[i].layout + routes[i].path) !== -1) {
	// 				return routes[i].path;
	// 			}
	// 		}
	// 	}
	// 	return activeNavbar;
	// };
	// const getRoutesnew = (
	// 	routes: RoutesType[]
	// ): any => {
	// 	return routes.map((prop, key) => {
	// 		if (prop.layout === '/auth') {
	// 			return <Route path={prop.layout + prop.path} component={prop.component} key={key} />;
	// 		}
	// 		if (prop.collapse && prop.items) {
	// 			return getRoutes(prop.items);
	// 		}

	// 		return null;
	// 	});
	// };
	const getRoutes = (routes: RoutesType[]): any => {
		return routes.map((prop, key) => {
			if (prop.layout === '/auth') {
				return <Route path={prop.layout + prop.path} component={prop.component} key={key} />;
			}
			if (prop.collapse) {
				return getRoutes(prop.items);
			}
			return null;
		});
	};
	const authBg = useColorModeValue('white', 'navy.900');
	document.documentElement.dir = 'ltr';
	return (
		<Box>
			<SidebarContext.Provider
				value={{
					toggleSidebar,
					setToggleSidebar
				}}>
				<Box
					bg={authBg}
					float='right'
					minHeight='100vh'
					height='100%'
					position='relative'
					w='100%'
					transition='all 0.33s cubic-bezier(0.685, 0.0473, 0.346, 1)'
					transitionDuration='.2s, .2s, .35s'
					transitionProperty='top, bottom, width'
					transitionTimingFunction='linear, linear, ease'>
					{getRoute() ? (
						<Box mx='auto' minH='100vh'>
							<Switch>
								{getRoutes(routes)}
								<Redirect from='/' to='/auth/sign-in/default' />
							</Switch>
						</Box>
					) : null}
				</Box>
			</SidebarContext.Provider>
		</Box>
	);
}
