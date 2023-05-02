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

// Chakra imports
import {
	Box,
	Button,
	Flex,
	FormControl,
	PinInput,
	PinInputField,
	Heading,
	Text,
	useColorModeValue
} from '@chakra-ui/react';
// Assets
import illustration from 'assets/img/auth/auth.png';
import DefaultAuth from 'layouts/auth/variants/Default';
import React from 'react';

function ForgotPassword() {
	// Chakra color mode
	const textColor = useColorModeValue('navy.700', 'white');
	const textColorDetails = useColorModeValue('navy.700', 'secondaryGray.600');
	const textColorBrand = useColorModeValue('brand.500', 'white');
	const borderColor = useColorModeValue('secondaryGray.400', 'whiteAlpha.100');
	return (
		<DefaultAuth illustrationBackground={illustration} image={illustration}>
			<Flex
				w='100%'
				maxW='max-content'
				mx={{ base: 'auto', lg: '0px' }}
				me='auto'
				h='100%'
				alignItems='start'
				justifyContent='center'
				mb={{ base: '30px', md: '60px', lg: '90px', xl: '60px' }}
				px={{ base: '25px', md: '0px' }}
				mt={{ base: '40px', lg: '14vh', xl: '22vh' }}
				flexDirection='column'>
				<Box me='auto' mb='34px'>
					<Heading
						color={textColor}
						fontSize='36px'
						mb='16px'
						mx={{ base: 'auto', lg: 'unset' }}
						textAlign={{ base: 'center', lg: 'left' }}>
						2-Step Verification
					</Heading>
					<Text
						color='gray.400'
						fontSize='md'
						maxW={{ base: '95%', md: '100%' }}
						mx={{ base: 'auto', lg: 'unset' }}
						textAlign={{ base: 'center', lg: 'left' }}>
						Enter your 2-Step Verification email code to unlock!
					</Text>
				</Box>
				<Flex
					zIndex='2'
					direction='column'
					w={{ base: '100%', md: '395px' }}
					maxW='100%'
					background='transparent'
					borderRadius='15px'
					mx={{ base: 'auto', lg: 'unset' }}
					me='auto'
					mb={{ base: '20px', md: 'auto' }}>
					<FormControl>
						<Flex justify='center'>
							<PinInput mx='auto' otp>
								<PinInputField
									fontSize='36px'
									color={textColor}
									borderRadius='16px'
									borderColor={borderColor}
									h={{ base: '63px', md: '95px' }}
									w={{ base: '63px', md: '95px' }}
									me='10px'
								/>
								<PinInputField
									fontSize='36px'
									color={textColor}
									borderRadius='16px'
									borderColor={borderColor}
									h={{ base: '63px', md: '95px' }}
									w={{ base: '63px', md: '95px' }}
									me='10px'
								/>
								<PinInputField
									fontSize='36px'
									color={textColor}
									borderRadius='16px'
									borderColor={borderColor}
									h={{ base: '63px', md: '95px' }}
									w={{ base: '63px', md: '95px' }}
									me='10px'
								/>
								<PinInputField
									fontSize='36px'
									color={textColor}
									borderRadius='16px'
									borderColor={borderColor}
									h={{ base: '63px', md: '95px' }}
									w={{ base: '63px', md: '95px' }}
								/>
							</PinInput>
						</Flex>

						<Button
							fontSize='14px'
							variant='brand'
							borderRadius='16px'
							fontWeight='500'
							w='100%'
							h='50'
							mb='24px'
							mt='12px'>
							Unlock
						</Button>
					</FormControl>
					<Flex flexDirection='column' justifyContent='center' alignItems='start' maxW='100%' mt='0px'>
						<Text
							color={textColorDetails}
							fontWeight='400'
							fontSize='14px'
							mx={{ base: 'auto', lg: 'unset' }}
							textAlign={{ base: 'center', lg: 'left' }}>
							Haven't received it?
							<Text color={textColorBrand} as='span' ms='5px' fontWeight='500'>
								Resend a new code
							</Text>
						</Text>
					</Flex>
				</Flex>
			</Flex>
		</DefaultAuth>
	);
}

export default ForgotPassword;
