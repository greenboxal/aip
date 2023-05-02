/* eslint-disable */

import { NavLink, useLocation } from 'react-router-dom';
// chakra imports
import { 
  Accordion,
  AccordionButton,
  AccordionIcon,
  AccordionItem,
  AccordionPanel,
  Box, 
  Flex, 
  HStack, 
  Text,
  List,
  Icon,
  ListItem,
  useColorModeValue 
} from '@chakra-ui/react';
// Assets
import { FaCircle } from "react-icons/fa";

export function SidebarLinks(props: {
	routes:RoutesType[];
	[x: string]: any;
}) {
	//   Chakra color mode
	let location = useLocation();
	let activeColor = useColorModeValue('gray.700', 'white');
	let inactiveColor = useColorModeValue('secondaryGray.600', 'secondaryGray.600');
	let activeIcon = useColorModeValue('brand.500', 'white'); 

	const { routes } = props;

	// verifies if routeName is the one active (in browser input)
	const activeRoute = (routeName: string) => {
		return location.pathname.includes(routeName);
	};
  
  // this function creates the links and collapses that appear in the sidebar (left menu)
  	const createLinks = (
		routes: RoutesType[]
	) => {
    return routes.map((route, key) => {
      if (route.collapse) {
        return (
          <Accordion allowToggle key={key}>
            <AccordionItem border='none' key={key}>
              <AccordionButton
                display='flex'
                alignItems='center'
                justifyContent='center'
                _hover={{
                  bg: "unset",
                }}
                _focus={{
                  boxShadow: "none",
                }}
                borderRadius='8px'
                w={{
                  sm: "100%",
                  xl: "100%",
                  "2xl": "95%",
                }}
                px={route.icon ? null : "0px"}
                py='0px'
                bg={"transparent"}
                ms={0}>
                {route.icon ? (
                  <Flex align='center' justifyContent='space-between' w='100%'>
                    <HStack
                      mb='6px'
                      spacing={
                        activeRoute(route.path.toLowerCase()) ? "22px" : "26px"
                      }>
                      <Flex
                        w='100%'
                        alignItems='center'
                        justifyContent='center'>
                        <Box
                          color={
                            activeRoute(route.path.toLowerCase())
                              ? activeIcon
                              : inactiveColor
                          }
                          me='12px'
                          mt='6px'>
                          {route.icon}
                        </Box>
                        <Text
                          me='auto'
                          color={
                            activeRoute(route.path.toLowerCase())
                              ? activeColor
                              : "secondaryGray.600"
                          }
                          fontWeight='500'
                          fontSize='md'>
                          {route.name}
                        </Text>
                      </Flex>
                    </HStack>
                    <AccordionIcon
                      ms='auto'
                      color={"secondaryGray.600"}
                      transform={route.icon ? null : "translateX(-70%)"}
                    />
                  </Flex>
                ) : (
                  <Flex pt='0px' pb='10px' alignItems='center' w='100%'>
                    <HStack
                      spacing={
                        activeRoute(route.path.toLowerCase()) ? "22px" : "26px"
                      }
                      ps='34px'>
                      <Text
                        me='auto'
                        color={
                          activeRoute(route.path.toLowerCase())
                            ? activeColor
                            : inactiveColor
                        }
                        fontWeight='500'
                        fontSize='sm'>
                        {route.name}
                      </Text>
                    </HStack>
                    <AccordionIcon
                      ms='auto'
                      color={"secondaryGray.600"}
                      transform={null}
                    />
                  </Flex>
                )}
              </AccordionButton>
              <AccordionPanel
                pe={route.icon ? null : "0px"}
                py='0px'
                ps={route.icon ? null : "8px"}>
                <List>
                  {
                    route.icon
                      ? createLinks(route.items) // for bullet accordion links
                      : createAccordionLinks(route.items) // for non-bullet accordion links
                  }
                </List>
              </AccordionPanel>
            </AccordionItem>
          </Accordion>
        );
      } else {
        return (
          <NavLink to={route.layout + route.path} key={key}>
            {route.icon ? (
              <Flex
                align='center'
                justifyContent='space-between'
                w='100%'
                ps='17px'
                mb='0px'>
                <HStack
                  mb='6px'
                  spacing={
                    activeRoute(route.path.toLowerCase()) ? "22px" : "26px"
                  }>
                  <Flex w='100%' alignItems='center' justifyContent='center'>
                    <Box
                      color={
                        activeRoute(route.path.toLowerCase())
                          ? activeIcon
                          : inactiveColor
                      }
                      me='12px'
                      mt='6px'>
                      {route.icon}
                    </Box>
                    <Text
                      me='auto'
                      color={
                        activeRoute(route.path.toLowerCase())
                          ? activeColor
                          : "secondaryGray.600"
                      }
                      fontWeight='500'>
                      {route.name}
                    </Text>
                  </Flex>
                </HStack>
              </Flex>
            ) : (
              <ListItem ms={null}>
                <Flex ps='34px' alignItems='center' mb='8px'>
                  <Text
                    color={
                      activeRoute(route.path.toLowerCase())
                        ? activeColor
                        : inactiveColor
                    }
                    fontWeight='500'
                    fontSize='sm'>
                    {route.name}
                  </Text>
                </Flex>
              </ListItem>
            )}
          </NavLink>
        );
      }
    });
  };
  // this function creates the links from the secondary accordions (for example auth -> sign-in -> default)
  const createAccordionLinks = (
		routes: RoutesType[]
	) => {
    return routes.map((route:RoutesType, key:number) => {
      return (
        <NavLink to={route.layout + route.path} key={key}>
          <ListItem
            ms='28px'
            display='flex'
            alignItems='center'
            mb='10px'
            key={key}>
            <Icon w='6px' h='6px' me='8px' as={FaCircle} color={activeIcon} />
            <Text
              color={
                activeRoute(route.path.toLowerCase())
                  ? activeColor
                  : inactiveColor
              }
              fontWeight={
                activeRoute(route.path.toLowerCase()) ? "bold" : "normal"
              }
              fontSize='sm'>
              {route.name}
            </Text>
          </ListItem>
        </NavLink>
      );
    });
  };
	//  BRAND
	return <>{createLinks(routes)}</>
}

export default SidebarLinks;
