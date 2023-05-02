import React from 'react';

// Chakra imports
import { Button, Flex, Link, Text } from '@chakra-ui/react';

// Assets
import banner from 'src/horizon/assets/img/product/OverviewBanner.png';

export default function Banner() {
	// Chakra Color Mode
	return (
		<Flex
			direction='column'
			bgImage={banner}
			bgSize='cover'
			mb='20px'
			py={{ base: '30px', md: '56px' }}
			px={{ base: '30px', md: '64px' }}
			borderRadius='30px'>
			<Text
				fontSize={{ base: '24px', md: '34px' }}
				color='white'
				mb='14px'
				maxW={{
					base: '75%',
					md: '54%',
					lg: '65%',
					xl: '65%',
					'2xl': '50%',
					'3xl': '42%'
				}}
				fontWeight='700'
				lineHeight={{ base: '32px', md: '42px' }}>
				Shop’s Products List Overview
			</Text>
			<Text
				fontSize='md'
				color='#E3DAFF'
				maxW={{
					base: '90%',
					md: '64%',
					lg: '66%',
					xl: '56%',
					'2xl': '46%',
					'3xl': '38%'
				}}
				fontWeight='500'
				mb='40px'
				lineHeight='28px'>
				Check your products statistics daily, see all sales, ratings, comments and updates here. You can edit
				and manage everything in the options panel.
			</Text>
			<Flex align='center'>
				<Button
					bg='white'
					color='black'
					_hover={{ bg: 'whiteAlpha.900' }}
					_active={{ bg: 'white' }}
					_focus={{ bg: 'white' }}
					fontWeight='500'
					fontSize='14px'
					py='20px'
					px='27'
					me='38px'>
					Sales Analytics
				</Button>
				<Link>
					<Text color='white' fontSize='sm' fontWeight='500'>
						Marketing Options
					</Text>
				</Link>
			</Flex>
		</Flex>
	);
}
