// Chakra imports
import { Box, Flex, Text, Badge, LightMode } from '@chakra-ui/react';
import LineChart from 'src/horizon/components/charts/LineChart';
// Custom components
import { lineChartDataSidebar, lineChartOptionsSidebar } from 'src/horizon/variables/charts';
export default function SidebarDocs() {
	const bgColor = 'linear-gradient(135deg, #868CFF 0%, #4318FF 100%)';

	return (
		<Flex
			justify='center'
			direction='column'
			align='center'
			bg={bgColor}
			borderRadius='30px'
			me='20px'
			position='relative'>
			<Flex direction='column' mb='12px' align='center' justify='center' px='15px' pt='30px'>
				<Text
					fontSize={{ base: 'lg', xl: '2xl' }}
					color='white'
					fontWeight='bold'
					lineHeight='150%'
					textAlign='center'
					px='10px'>
					$3942.58
				</Text>
				<Text fontSize='sm' color='white' px='10px' mb='14px' textAlign='center'>
					Total balance
				</Text>
				<LightMode>
					<Badge colorScheme='green' color='green.500' size='lg' borderRadius='58px'>
						+2.45%
					</Badge>
				</LightMode>
				<Box h='160px'>
					<LineChart chartData={lineChartDataSidebar} chartOptions={lineChartOptionsSidebar} />
				</Box>
			</Flex>
		</Flex>
	);
}
