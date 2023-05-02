// Chakra imports
import { Box, Flex, Text } from '@chakra-ui/react';
import Footer from 'src/horizon/components/footer/FooterAuthCentered';
// Custom components
import FixedPlugin from 'src/horizon/components/fixedPlugin/FixedPlugin';
import Navbar from 'src/horizon/components/navbar/NavbarAuth';

function AuthPricing(props: {
	children: JSX.Element;
	title?: string;
	description?: string;
	image?: string;
	contentTop?: string | number | { [x: string]: string | number };
	contentBottom?: string | number | { [x: string]: string | number };
}) {
	const { children, title, description, contentTop, contentBottom } = props;
	return (
		<Flex
			direction='column'
			alignSelf='center'
			justifySelf='center'
			overflow='hidden'
			mx={{ base: '10px', lg: '0px' }}
			minH='100vh'>
			<Box
				position='absolute'
				minH={{ base: '60vh', md: '60vh' }}
				maxH={{ base: '60vh', md: '60vh' }}
				w={{ md: 'calc(100vw)' }}
				maxW={{ md: 'calc(100vw)' }}
				left='0'
				right='0'
				bgRepeat='no-repeat'
				overflow='hidden'
				zIndex='0'
				top='0'
				bgImage={'linear-gradient(135deg, #868CFF 0%, #4318FF 100%)'}
				mx={{ md: 'auto' }}
			/>
			<Navbar />
			<Flex
				w={{ base: '100%', md: 'max-content' }}
				p={{ base: '10px', md: '50px' }}
				h='max-content'
				mx='auto'
				maxW='100%'
				mt={contentTop}
				mb={contentBottom}>
				{title && description ? (
					<Flex
						direction='column'
						textAlign='center'
						justifyContent='center'
						align='center'
						mt='125px'
						mb='30px'>
						<Text fontSize='4xl' color='white' fontWeight='bold'>
							{title}
						</Text>
						<Text
							fontSize='md'
							color='white'
							fontWeight='normal'
							mt='10px'
							mb='26px'
							w={{ base: '90%', sm: '60%', lg: '40%', xl: '333px' }}>
							{description}
						</Text>
					</Flex>
				) : null}
				{children}
			</Flex>
			<Footer />
			<FixedPlugin />
		</Flex>
	);
}
// PROPS

export default AuthPricing;
