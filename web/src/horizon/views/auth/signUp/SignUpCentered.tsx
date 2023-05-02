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

import { NavLink } from 'react-router-dom';
import React from 'react';

// Chakra imports
import {
	Box,
	Button,
	Checkbox,
	Flex,
	FormControl,
	FormLabel,
	Heading,
	Icon,
	Input,
	InputGroup,
	InputRightElement,
	Link,
	SimpleGrid,
	Text,
	useColorModeValue
} from '@chakra-ui/react';

// Custom components
import { HSeparator } from 'src/horizon/components/separator/Separator';
import CenteredAuth from 'src/horizon/layouts/auth/variants/Centered';

// Assets
import { FcGoogle } from 'react-icons/fc';
import { MdOutlineRemoveRedEye } from 'react-icons/md';
import { RiEyeCloseLine } from 'react-icons/ri';

function SignUp() {
	// Chakra color mode
	const textColor = useColorModeValue('navy.700', 'white');
	const textColorSecondary = 'gray.400';
	const textColorDetails = useColorModeValue('navy.700', 'secondaryGray.600');
	const textColorBrand = useColorModeValue('brand.500', 'white');
	const brandStars = useColorModeValue('brand.500', 'brand.400');
	const googleBg = useColorModeValue('secondaryGray.300', 'whiteAlpha.200');
	const googleText = useColorModeValue('navy.700', 'white');
	const googleHover = useColorModeValue({ bg: 'gray.200' }, { bg: 'whiteAlpha.300' });
	const googleActive = useColorModeValue({ bg: 'secondaryGray.300' }, { bg: 'whiteAlpha.200' });
	const [ show, setShow ] = React.useState(false);
	const handleClick = () => setShow(!show);
	return (
		<CenteredAuth
			image={'linear-gradient(135deg, #868CFF 0%, #4318FF 100%)'}
			cardTop={{ base: '140px', md: '14vh' }}
			cardBottom={{ base: '50px', lg: '100px' }}>
			<Flex
				maxW='max-content'
				mx={{ base: 'auto', lg: '0px' }}
				me='auto'
				justifyContent='center'
				px={{ base: '20px', md: '0px' }}
				flexDirection='column'>
				<Box me='auto'>
					<Heading color={textColor} fontSize={{ base: '34px', lg: '36px' }} mb='10px'>
						Sign Up
					</Heading>
					<Text mb='36px' ms='4px' color={textColorSecondary} fontWeight='400' fontSize='md'>
						Enter your email and password to sign up!
					</Text>
				</Box>
				<Flex
					zIndex='2'
					direction='column'
					w={{ base: '100%', md: '420px' }}
					maxW='100%'
					background='transparent'
					borderRadius='15px'
					mx={{ base: 'auto', lg: 'unset' }}
					me='auto'
					mb={{ base: '20px', md: 'auto' }}>
					<Button
						fontSize='sm'
						me='0px'
						mb='26px'
						py='15px'
						h='50px'
						borderRadius='16px'
						bg={googleBg}
						color={googleText}
						fontWeight='500'
						_hover={googleHover}
						_active={googleActive}
						_focus={googleActive}>
						<Icon as={FcGoogle} w='20px' h='20px' me='10px' />
						Sign up with Google
					</Button>
					<Flex align='center' mb='25px'>
						<HSeparator />
						<Text color={textColorSecondary} mx='14px'>
							or
						</Text>
						<HSeparator />
					</Flex>
					<FormControl>
						<SimpleGrid columns={{ base: 1, md: 2 }} gap={{ sm: '10px', md: '26px' }}>
							<Flex direction='column'>
								<FormLabel
									display='flex'
									ms='4px'
									fontSize='sm'
									fontWeight='500'
									color={textColor}
									mb='8px'>
									First name<Text color={brandStars}>*</Text>
								</FormLabel>
								<Input
									isRequired={true}
									fontSize='sm'
									ms={{ base: '0px', md: '4px' }}
									placeholder='First name'
									variant='auth'
									mb='24px'
									size='lg'
								/>
							</Flex>
							<Flex direction='column'>
								<FormLabel
									display='flex'
									ms='4px'
									fontSize='sm'
									fontWeight='500'
									color={textColor}
									mb='8px'>
									Last name<Text color={brandStars}>*</Text>
								</FormLabel>
								<Input
									isRequired={true}
									variant='auth'
									fontSize='sm'
									placeholder='Last name'
									mb='24px'
									size='lg'
								/>
							</Flex>
						</SimpleGrid>
						<FormLabel display='flex' ms='4px' fontSize='sm' fontWeight='500' color={textColor} mb='8px'>
							Email<Text color={brandStars}>*</Text>
						</FormLabel>
						<Input
							isRequired={true}
							variant='auth'
							fontSize='sm'
							type='email'
							placeholder='mail@simmmple.com'
							mb='24px'
							size='lg'
						/>
						<FormLabel ms='4px' fontSize='sm' fontWeight='500' color={textColor} display='flex'>
							Password<Text color={brandStars}>*</Text>
						</FormLabel>
						<InputGroup size='md'>
							<Input
								isRequired={true}
								variant='auth'
								fontSize='sm'
								ms={{ base: '0px', md: '4px' }}
								placeholder='Min. 8 characters'
								mb='24px'
								size='lg'
								type={show ? 'text' : 'password'}
							/>
							<InputRightElement display='flex' alignItems='center' mt='4px'>
								<Icon
									color={textColorSecondary}
									_hover={{ cursor: 'pointer' }}
									as={show ? RiEyeCloseLine : MdOutlineRemoveRedEye}
									onClick={handleClick}
								/>
							</InputRightElement>
						</InputGroup>
						<Flex justifyContent='space-between' align='center' mb='24px'>
							<FormControl display='flex' alignItems='start'>
								<Checkbox id='remember-login' colorScheme='brandScheme' me='10px' mt='3px' />
								<FormLabel
									htmlFor='remember-login'
									mb='0'
									fontWeight='normal'
									color={textColor}
									fontSize='sm'>
									By creating an account means you agree to the{' '}
									<Link href='https://simmmple.com/terms-of-service' fontWeight='500'>
										Terms and Conditions,
									</Link>{' '}
									and our{' '}
									<Link href='https://simmmple.com/privacy-policy' fontWeight='500'>
										Privacy Policy
									</Link>
								</FormLabel>
							</FormControl>
						</Flex>
						<Button variant='brand' fontSize='14px' fontWeight='500' w='100%' h='50' mb='24px'>
							Create my account
						</Button>
					</FormControl>
					<Flex flexDirection='column' justifyContent='center' alignItems='start' maxW='100%' mt='0px'>
						<Text color={textColorDetails} fontWeight='400' fontSize='sm'>
							Already a member?
							<NavLink to='/auth/sign-in'>
								<Text color={textColorBrand} as='span' ms='5px' fontWeight='500'>
									Sign in
								</Text>
							</NavLink>
						</Text>
					</Flex>
				</Flex>
			</Flex>
		</CenteredAuth>
	);
}

export default SignUp;
