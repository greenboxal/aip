/*!
  _   _  ___  ____  ___ ________  _   _   _   _ ___   ____  ____   ___  
 | | | |/ _ \|  _ \|_ _|__  / _ \| \ | | | | | |_ _| |  _ \|  _ \ / _ \ 
 | |_| | | | | |_) || |  / / | | |  \| | | | | || |  | |_) | |_) | | | |
 |  _  | |_| |  _ < | | / /| |_| | |\  | | |_| || |  |  __/|  _ <| |_| |
 |_| |_|\___/|_| \_\___/____\___/|_| \_|  \___/|___| |_|   |_| \_\\___/ 
                                                                                                                                                                                                                                                                                                                                       
=========================================================
* Horizon UI Dashboard PRO - v1.0.0
=========================================================

* Product Page: https://www.horizon-ui.com/pro/
* Copyright 2022 Horizon UI (https://www.horizon-ui.com/)

* Designed and Coded by Simmmple

=========================================================

* The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

*/

// Since this is a specific page, if you want to make it work like a real chat, we have some suggestions:
// 1. Add to the `Messages` components a `data` prop which will pass the messages fetched from a server.
// 2. Put all `ChatHeader` as Chakra UI `Tab` and each of them will open a panel, which will be your `Messages` component with its data.

// These are purely suggestions, although in Horizon v2.0 we will make an example with a simulated chat !

// Chakra imports
import {
	Box,
	Button,
	Drawer,
	DrawerBody,
	DrawerContent,
	Flex,
	Grid,
	Icon,
	IconButton,
	Input,
	InputGroup,
	InputLeftElement,
	Menu,
	MenuButton,
	MenuItem,
	MenuList,
	Text,
	DrawerCloseButton,
	useColorModeValue,
	useDisclosure
} from '@chakra-ui/react';
// Assets
import avatar1 from 'src/horizon/assets/img/avatars/avatar1.png';
import avatar2 from 'src/horizon/assets/img/avatars/avatar2.png';
import avatar4 from 'src/horizon/assets/img/avatars/avatar4.png';
import avatar5 from 'src/horizon/assets/img/avatars/avatar5.png';
import avatar6 from 'src/horizon/assets/img/avatars/avatar6.png';
import avatar7 from 'src/horizon/assets/img/avatars/avatar7.png';
import avatar8 from 'src/horizon/assets/img/avatars/avatar8.png';
import avatar9 from 'src/horizon/assets/img/avatars/avatar9.png';
import { FiSearch } from 'react-icons/fi';
import { MdChevronLeft } from 'react-icons/md';
import { FaRegEdit } from 'react-icons/fa';
import {
	MdOutlineCardTravel,
	MdOutlineLightbulb,
	MdOutlineMoreVert,
	MdOutlinePerson,
	MdOutlineSettings
} from 'react-icons/md';
// Custom components
import Card from 'src/horizon/components/card/Card';
import Messages from 'src/horizon/views/admin/main/others/messages/components/Messages';
import ChatHeader from 'src/horizon/components/chat/ChatHeader';
// import Conversations from 'views/admin/main/others/messages/components/Conversations';

export default function Default() {
	// Chakra Color Mode
	const textColor = useColorModeValue('secondaryGray.900', 'white');
	const searchIconColor = useColorModeValue('gray.700', 'white');
	const inputText = useColorModeValue('gray.700', 'gray.100');
	const blockBg = useColorModeValue('secondaryGray.300', 'navy.700');
	const brandButton = useColorModeValue('brand.500', 'brand.400');
	// Ellipsis modals
	const { isOpen: isOpen1, onOpen: onOpen1, onClose: onClose1 } = useDisclosure();

	// Chakra Color Mode
	const textHover = useColorModeValue(
		{ color: 'secondaryGray.900', bg: 'unset' },
		{ color: 'secondaryGray.500', bg: 'unset' }
	);
	const bgList = useColorModeValue('white', 'whiteAlpha.100');
	const bgShadow = useColorModeValue('14px 17px 40px 4px rgba(112, 144, 176, 0.08)', 'unset');
	let sidebarBackgroundColor = useColorModeValue('white', 'navy.800');
	const { isOpen, onOpen, onClose } = useDisclosure();
	// Chakra Color Mode
	return (
		<Grid
			display={{ base: 'block', md: 'grid' }}
			pt={{ base: '130px', md: '80px', xl: '80px' }}
			gridTemplateColumns={{ base: '1.5fr 1.83fr', lg: '1fr 1.83fr' }}
			gap={{ base: '20px', xl: '20px' }}>
			<Card gridArea='1 / 1 / 2 / 2' p='0px' px='20px' pt='20px' pb='40px'>
				{/* <Conversations onOpenMessages={onOpen} /> */}
				<Box>
					<Box>
						<Flex mb='15px' align='center' justify='space-between'>
							<Text color={textColor} fontSize='xl' fontWeight='700'>
								Your Chats
							</Text>
							<Menu isOpen={isOpen1} onClose={onClose1}>
								<MenuButton onClick={onOpen1} mb='0px'>
									<Icon
										mb='-6px'
										cursor='pointer'
										as={MdOutlineMoreVert}
										color={textColor}
										maxW='min-content'
										maxH='min-content'
										w='24px'
										h='24px'
									/>
								</MenuButton>
								<MenuList
									w='150px'
									minW='unset'
									maxW='150px !important'
									border='transparent'
									backdropFilter='blur(63px)'
									bg={bgList}
									boxShadow={bgShadow}
									borderRadius='20px'
									p='15px'>
									<MenuItem
										transition='0.2s linear'
										color={textColor}
										_hover={textHover}
										p='0px'
										borderRadius='8px'
										_active={{
											bg: 'transparent'
										}}
										_focus={{
											bg: 'transparent'
										}}
										mb='10px'>
										<Flex align='center'>
											<Icon as={MdOutlinePerson} h='16px' w='16px' me='8px' />
											<Text fontSize='sm' fontWeight='400'>
												Panel 1
											</Text>
										</Flex>
									</MenuItem>
									<MenuItem
										transition='0.2s linear'
										p='0px'
										borderRadius='8px'
										color={textColor}
										_hover={textHover}
										_active={{
											bg: 'transparent'
										}}
										_focus={{
											bg: 'transparent'
										}}
										mb='10px'>
										<Flex align='center'>
											<Icon as={MdOutlineCardTravel} h='16px' w='16px' me='8px' />
											<Text fontSize='sm' fontWeight='400'>
												Panel 2
											</Text>
										</Flex>
									</MenuItem>
									<MenuItem
										transition='0.2s linear'
										p='0px'
										borderRadius='8px'
										color={textColor}
										_hover={textHover}
										_active={{
											bg: 'transparent'
										}}
										_focus={{
											bg: 'transparent'
										}}
										mb='10px'>
										<Flex align='center'>
											<Icon as={MdOutlineLightbulb} h='16px' w='16px' me='8px' />
											<Text fontSize='sm' fontWeight='400'>
												Panel 3
											</Text>
										</Flex>
									</MenuItem>
									<MenuItem
										transition='0.2s linear'
										color={textColor}
										_hover={textHover}
										p='0px'
										borderRadius='8px'
										_active={{
											bg: 'transparent'
										}}
										_focus={{
											bg: 'transparent'
										}}>
										<Flex align='center'>
											<Icon as={MdOutlineSettings} h='16px' w='16px' me='8px' />
											<Text fontSize='sm' fontWeight='400'>
												Panel 4
											</Text>
										</Flex>
									</MenuItem>
								</MenuList>
							</Menu>
						</Flex>
						<Flex align='center' w='calc(100%)' bottom='20px'>
							<InputGroup me='10px' w={{ base: '100%' }}>
								<InputLeftElement
									zIndex='1'
									children={
										<IconButton
											aria-label='iconbutton'
											bg='inherit'
											borderRadius='inherit'
											_hover={{ bg: 'none' }}
											_active={{
												bg: 'inherit',
												transform: 'none',
												borderColor: 'transparent'
											}}
											_focus={{
												boxShadow: 'none'
											}}
											icon={<Icon as={FiSearch} color={searchIconColor} w='15px' h='15px' />}
										/>
									}
								/>
								<Input
									variant='search'
									fontSize='sm'
									pl='35px !important'
									h='40px'
									bg={blockBg}
									color={inputText}
									fontWeight='500'
									_placeholder={{ color: 'gray.400', fontSize: '14px' }}
									borderRadius={'50px'}
									placeholder={'Search'}
								/>
							</InputGroup>
							<Button
								borderRadius='50%'
								ms={{ base: '14px', md: 'auto' }}
								bg={brandButton}
								w={{ base: '35px', md: '40px' }}
								h={{ base: '35px', md: '40px' }}
								minW={{ base: '35px', md: '40px' }}
								minH={{ base: '35px', md: '40px' }}
								variant='no-hover'>
								<Icon
									as={FaRegEdit}
									color='white'
									w={{ base: '16px', md: '16px' }}
									h={{ base: '16px', md: '16px' }}
								/>
							</Button>
						</Flex>
					</Box>
					<Box onClick={window.innerWidth < 450 ? onOpen : null}>
						<ChatHeader
							name='Roberto Michael'
							lastMessage='Hi there, How are you? All good?'
							sum='-$15.50'
							avatar={avatar2}
							hour='09:00 PM'
						/>
					</Box>
					<Box onClick={window.innerWidth < 450 ? onOpen : null}>
						<ChatHeader
							name='Emily James'
							lastMessage='Be careful, it’s raining outside! :)'
							sum='-$15.50'
							avatar={avatar1}
							hour='08:45 PM'
						/>
					</Box>
					<Box onClick={window.innerWidth < 450 ? onOpen : null}>
						<ChatHeader
							name='Alexander Parker'
							lastMessage='It contains a lot of good lessons about effective...'
							sum='-$15.50'
							avatar={avatar5}
							hour='08:42 PM'
						/>
					</Box>
					<Box onClick={window.innerWidth < 450 ? onOpen : null}>
						<ChatHeader
							name='Esthera William'
							lastMessage='Wow! This picture is amazing! Send me more!'
							sum='-$15.50'
							avatar={avatar4}
							hour='06:32 PM'
						/>
					</Box>
					<Box onClick={window.innerWidth < 450 ? onOpen : null}>
						<ChatHeader
							name='Lawrence Peter'
							lastMessage='You look so amazing today!'
							sum='-$15.50'
							avatar={avatar8}
							hour='06:30 PM'
						/>
					</Box>
					<Box onClick={window.innerWidth < 450 ? onOpen : null}>
						<ChatHeader
							name='Iaon Dint'
							lastMessage='I’m back from Belgium, do you want to meet?'
							sum='-$15.50'
							avatar={avatar9}
							hour='05:57 PM'
						/>
					</Box>
					<Box onClick={window.innerWidth < 450 ? onOpen : null}>
						<ChatHeader
							name='William Jackson'
							lastMessage='That’s awesome!!! What technology do you used...'
							sum='-$15.50'
							avatar={avatar6}
							hour='04:32 PM'
						/>
					</Box>
					<Box onClick={window.innerWidth < 450 ? onOpen : null}>
						<ChatHeader
							last
							name='Markus Aurelius'
							lastMessage='Hello! How you doin’? I’m going to Italy this week...'
							sum='-$15.50'
							avatar={avatar7}
							hour='01:08 PM'
						/>
					</Box>
				</Box>
			</Card>
			<Card display={{ base: 'none', md: 'flex' }} h='100%' gridArea='1 / 2 / 2 / 3' px='0px'>
				<Messages status='active' name='Roberto Michael' />
			</Card>
			<Drawer isOpen={isOpen} placement='left' onClose={onClose}>
				<DrawerContent
					w='100%'
					maxW='calc(100vw - 32px)'
					ms={{
						sm: '16px'
					}}
					my={{
						sm: '16px'
					}}
					borderRadius='16px'
					bg={sidebarBackgroundColor}>
					<DrawerCloseButton
						boxShadow='unset !important'
						_hover={{ bg: 'transparent' }}
						left='0px !important'
						top='18px'
						color='secondaryGray.700'>
						<Icon as={MdChevronLeft} color={searchIconColor} w='28px' h='28px' />
					</DrawerCloseButton>
					<DrawerBody px='0px'>
						<Messages status='active' name='Roberto Michael' />
					</DrawerBody>
				</DrawerContent>
			</Drawer>
		</Grid>
	);
}
