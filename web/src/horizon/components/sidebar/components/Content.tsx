// chakra imports
import { Avatar, Box, Flex, Stack, Text, useColorModeValue } from '@chakra-ui/react';
//   Custom components
import Brand from 'src/horizon/components/sidebar/components/Brand';
import Links from 'src/horizon/components/sidebar/components/Links';
import SidebarCard from 'src/horizon/components/sidebar/components/SidebarCard';
import avatar4 from 'src/horizon/assets/img/avatars/avatar4.png';

// FUNCTIONS

function SidebarContent(props: { routes: RoutesType[] }) {
	const { routes } = props;
	const textColor = useColorModeValue('navy.700', 'white');
	// SIDEBAR
	return (
		<Flex direction='column' height='100%' pt='25px' borderRadius='30px'>
			<Brand />
			<Stack direction='column' mb='auto' mt='8px'>
				<Box ps='20px' pe={{ md: '16px', '2xl': '1px' }}>
					<Links routes={routes} />
				</Box>
			</Stack>

			<Box ps='20px' pe={{ md: '16px', '2xl': '0px' }} mt='60px' borderRadius='30px'>
				<SidebarCard />
			</Box>
			<Flex mt='75px' mb='56px' justifyContent='center' alignItems='center'>
				<Avatar h='48px' w='48px' src={avatar4} me='20px' />
				<Box>
					<Text color={textColor} fontSize='md' fontWeight='700'>
						Adela Parkson
					</Text>
					<Text color='secondaryGray.600' fontSize='sm' fontWeight='400'>
						Product Designer
					</Text>
				</Box>
			</Flex>
		</Flex>
	);
}

export default SidebarContent;
