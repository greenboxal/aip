// Chakra imports
import {
	Avatar,
	Box,
	Button,
	Flex,
	Icon,
	IconButton,
	Input,
	InputGroup,
	InputLeftElement,
	InputRightElement,
	Link,
	Menu,
	MenuButton,
	MenuItem,
	MenuList,
	Text,
	useColorModeValue,
	useDisclosure
} from '@chakra-ui/react';
import { Scrollbars } from 'react-custom-scrollbars-2';
import MessageBlock from 'src/horizon/components/chat/MessageBlock';
import React from 'react';
// Assets
import { messagesRenderThumb, messagesRenderTrack, messagesRenderView } from 'src/horizon/components/scrollbar/Scrollbar';
import { FaCircle } from 'react-icons/fa';
import { FiSearch } from 'react-icons/fi';
import { IoPaperPlane } from 'react-icons/io5';
import {
	MdTagFaces,
	MdOutlineCardTravel,
	MdOutlineLightbulb,
	MdOutlineMoreVert,
	MdOutlinePerson,
	MdOutlineSettings,
	MdOutlineImage,
	MdAttachFile,
	MdAdd
} from 'react-icons/md';
import avatar2 from 'src/horizon/assets/img/avatars/avatar2.png';

export default function Messages(props: { status: string; name: string; [x: string]: any }) {
	const { status, name, ...rest } = props;

	// Chakra Color Mode
	const textColor = useColorModeValue('secondaryGray.900', 'white');

	const inputColor = useColorModeValue('secondaryGray.700', 'secondaryGray.700');
	const inputText = useColorModeValue('gray.700', 'gray.100');
	const blockBg = useColorModeValue('secondaryGray.300', 'navy.700');
	const brandButton = useColorModeValue('brand.500', 'brand.400');
	const bgInput = useColorModeValue(
		'linear-gradient(1.02deg, #FFFFFF 49.52%, rgba(255, 255, 255, 0) 99.07%)',
		'linear-gradient(1.02deg, #111C44 49.52%, rgba(17, 28, 68, 0) 99.07%)'
	);
	// Ellipsis modals
	const { isOpen: isOpen1, onOpen: onOpen1, onClose: onClose1 } = useDisclosure();

	// Chakra Color Mode
	const textHover = useColorModeValue(
		{ color: 'secondaryGray.900', bg: 'unset' },
		{ color: 'secondaryGray.500', bg: 'unset' }
	);
	const bgList = useColorModeValue('white', 'whiteAlpha.100');
	const brandColor = useColorModeValue('brand.500', 'white');
	const bgShadow = useColorModeValue('14px 17px 40px 4px rgba(112, 144, 176, 0.08)', 'unset');
	const borderColor = useColorModeValue('secondaryGray.400', 'whiteAlpha.100');
	return (
		<Box h='100%' {...rest}>
			<Flex px='34px' pb='25px' borderBottom='1px solid' borderColor={borderColor} align='center'>
				<Avatar
					h={{ base: '40px', '2xl': '50px' }}
					w={{ base: '40px', '2xl': '50px' }}
					src={avatar2}
					me='16px'
				/>
				<Box>
					<Text color={textColor} fontSize={{ base: 'md', md: 'xl' }} fontWeight='700'>
						Roberto Michael
					</Text>
					<Flex align='center'>
						<Icon
							w='6px'
							h='6px'
							me='8px'
							as={FaCircle}
							color={status === 'active' ? 'green.500' : status === 'away' ? 'orange.500' : 'red.500'}
						/>
						<Text fontSize={{ base: 'sm', md: 'md' }}>
							{status === 'active' ? 'Active' : status === 'away' ? 'Away' : 'Offline'}{' '}
						</Text>
					</Flex>
				</Box>
				<Flex align='center' ms='auto' />
				<Menu isOpen={isOpen1} onClose={onClose1}>
					<MenuButton onClick={onOpen1} mb='0px' me='8px'>
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
					</MenuList>{' '}
				</Menu>
				<Icon cursor='pointer' as={FiSearch} color={textColor} w='24px' h='24px' />
			</Flex>
			<Box h='calc(100% - 80px)' px={{ base: '10px', md: '20px' }} pt='45px' position='relative'>
				<Scrollbars
					autoHide
					renderTrackVertical={messagesRenderTrack}
					renderThumbVertical={messagesRenderThumb}
					renderView={messagesRenderView}>
					<Flex overflow='hidden'>
						<Flex
							direction='column'
							w='100%'
							maxW={{ base: '90%', lg: 'calc(100% - 80px)' }}
							boxSizing='border-box'>
							<MessageBlock content='Hi there, How are you? All good?' time='09:00 PM' side='left' />
							<MessageBlock
								content='I saw an amazing dashboard called Horizon UI Dashboard, is made by Simmmple, I want to know what you think about it, because I like it so much! 😁'
								time='09:00 PM'
								side='left'
							/>
							<MessageBlock
								content={
									<span>
										Go and check it out! Here is the link:{' '} <br/>
										<Link
											color={brandColor}
											target='_blank'
											href='https://horizon-ui.com/chakra-pro/?ref=comments-page'>
										    horizon-ui.com/chakra-pro/
										</Link>
									</span>
								}
								time='09:00 PM'
								side='left'
							/>
						</Flex>
					</Flex>
					<Flex mb='50px' overflow='hidden' w='94%' ms='auto' justify='end'>
						<Flex
							direction='column'
							w='calc(90%)'
							maxW={{ base: '90%', lg: 'calc(100% - 80px)' }}
							boxSizing='border-box'
							alignItems='flex-end'>
							<MessageBlock
								seen
								content={
									<span>
										Hello, Roberto! Hope you are fine! Let me take a look! Sounds interesting!
									</span>
								}
								time='09:23 PM'
							/>
							<MessageBlock
								isLast
								seen
								content={
									<span>
										OMG!! It’s so innovative and awesome! I think I am going to buy it for my
										projects! It’s a game changer!!🔥
									</span>
								}
								time='09:25 PM'
							/>
						</Flex>
					</Flex>
				</Scrollbars>
				<Flex
					bg={bgInput}
					backdropFilter='blur(20px)'
					pt='10px'
					position='absolute'
					w={{ base: 'calc(100% - 20px)', md: 'calc(100% - 40px)' }}
					bottom='0px'>
					<InputGroup me='10px' w={{ base: '100%' }}>
						<InputLeftElement
							children={
								<Box>
									<IconButton
										aria-label='iconbutton'
										ms='25px'
										h='max-content'
										w='max-content'
										mt='28px'
										bg='inherit'
										borderRadius='inherit'
										display={{ base: 'none', lg: 'flex' }}
										_hover={{ bg: 'none' }}
										_active={{
											bg: 'inherit',
											transform: 'none',
											borderColor: 'transparent'
										}}
										_focus={{
											boxShadow: 'none'
										}}
										icon={<Icon as={MdTagFaces} color={inputColor} w='30px' h='30px' />}
									/>
									<IconButton
										aria-label='iconbutton'
										mt='10px'
										w='30px'
										h='30px'
										bg='inherit'
										borderRadius='inherit'
										display={{ base: 'flex', lg: 'none' }}
										_hover={{ bg: 'none' }}
										_active={{
											bg: 'inherit',
											transform: 'none',
											borderColor: 'transparent'
										}}
										_focus={{
											boxShadow: 'none'
										}}
										icon={<Icon as={MdAdd} color={inputColor} w='30px' h='30px' />}
									/>
								</Box>
							}
						/>
						<InputRightElement
							display={{ base: 'none', lg: 'flex' }}
							children={
								<Flex me='70px'>
									<IconButton
										aria-label='iconbutton'
										h='max-content'
										w='max-content'
										mt='28px'
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
										icon={<Icon as={MdAttachFile} color={inputColor} w='30px' h='30px' />}
									/>
									<IconButton
										aria-label='iconbutton'
										h='max-content'
										w='max-content'
										mt='28px'
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
										icon={<Icon as={MdOutlineImage} color={inputColor} w='30px' h='30px' />}
									/>
								</Flex>
							}
						/>
						<Input
							variant='search'
							fontSize='md'
							pl={{ base: '40px !important', lg: '65px !important' }}
							pr={{
								base: '0px',
								lg: '145px !important'
							}}
							h={{ base: '50px', lg: '70px' }}
							bg={blockBg}
							color={inputText}
							fontWeight='500'
							_placeholder={{ color: 'gray.400', fontSize: '16px' }}
							borderRadius={'50px'}
							placeholder={'Write your message...'}
						/>
					</InputGroup>
					<Button
						borderRadius='50%'
						ms={{ base: '14px', lg: 'auto' }}
						bg={brandButton}
						w={{ base: '50px', lg: '70px' }}
						h={{ base: '50px', lg: '70px' }}
						minW={{ base: '50px', lg: '70px' }}
						minH={{ base: '50px', lg: '70px' }}
						variant='no-hover'>
						<Icon
							as={IoPaperPlane}
							color='white'
							w={{ base: '18px', lg: '25px' }}
							h={{ base: '18px', lg: '25px' }}
						/>
					</Button>
				</Flex>
			</Box>
		</Box>
	);
}
