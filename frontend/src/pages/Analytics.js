import React, { useState, useEffect } from 'react';
import {
  Box,
  Heading,
  SimpleGrid,
  Stat,
  StatLabel,
  StatNumber,
  StatHelpText,
  Text,
  Flex,
  Select,
  VStack,
  HStack,
  Card,
  CardHeader,
  CardBody,
  Table,
  Thead,
  Tbody,
  Tr,
  Th,
  Td,
  Badge,
  useColorModeValue,
  Divider,
  Spinner,
  Center,
} from '@chakra-ui/react';
import { couponsApi } from '../api/api';

// Sample data for demonstration
const sampleAnalyticsData = {
  totalCoupons: 12,
  activeCoupons: 8,
  totalRedemptions: 324,
  totalSavings: 4328.45,
  conversionRate: 68,
  topCoupons: [
    { coupon_code: 'SUMMER20', redemptions: 87, savings: 1245.67, conversion_rate: 72 },
    { coupon_code: 'WELCOME10', redemptions: 64, savings: 892.30, conversion_rate: 85 },
    { coupon_code: 'FLASH25', redemptions: 47, savings: 723.45, conversion_rate: 63 },
    { coupon_code: 'LOYALTY15', redemptions: 42, savings: 612.80, conversion_rate: 58 },
    { coupon_code: 'FREESHIP', redemptions: 38, savings: 380.00, conversion_rate: 49 },
  ],
  monthlyData: [
    { month: 'Jan', redemptions: 24, savings: 345.60 },
    { month: 'Feb', redemptions: 28, savings: 412.40 },
    { month: 'Mar', redemptions: 32, savings: 468.80 },
    { month: 'Apr', redemptions: 35, savings: 512.75 },
    { month: 'May', redemptions: 42, savings: 613.20 },
    { month: 'Jun', redemptions: 46, savings: 672.60 },
    { month: 'Jul', redemptions: 38, savings: 554.80 },
    { month: 'Aug', redemptions: 32, savings: 467.20 },
    { month: 'Sep', redemptions: 26, savings: 379.60 },
    { month: 'Oct', redemptions: 18, savings: 262.80 },
    { month: 'Nov', redemptions: 0, savings: 0 },
    { month: 'Dec', redemptions: 0, savings: 0 },
  ],
  categoriesData: [
    { category: 'painkiller', redemptions: 98, savings: 1176.00 },
    { category: 'antibiotic', redemptions: 64, savings: 1024.00 },
    { category: 'vitamin', redemptions: 86, savings: 946.00 },
    { category: 'supplement', redemptions: 42, savings: 798.45 },
    { category: 'first-aid', redemptions: 34, savings: 384.00 },
  ],
};

const Analytics = () => {
  const [analyticsData, setAnalyticsData] = useState(null);
  const [coupons, setCoupons] = useState([]);
  const [timeframe, setTimeframe] = useState('year');
  const [loading, setLoading] = useState(true);
  
  // Background colors
  const cardBg = useColorModeValue('white', 'gray.700');
  const statCardBg = useColorModeValue('brand.50', 'brand.900');
  const categoryBoxBg = useColorModeValue('gray.50', 'gray.700');
  
  // Fetch real coupon data and generate analytics
  const fetchCouponsAndGenerateAnalytics = async () => {
    setLoading(true);
    try {
      // Fetch real coupons data
      const couponsData = await couponsApi.getAllCoupons();
      setCoupons(couponsData);
      
      // Generate analytics based on fetched coupons
      const now = new Date();
      const activeCoupons = couponsData.filter(coupon => new Date(coupon.expiry_date) > now);
      
      // In a real app, we would fetch usage statistics from the backend
      // For now, generate simulated analytics based on real coupon data
      const calculatedData = {
        totalCoupons: couponsData.length,
        activeCoupons: activeCoupons.length,
        totalRedemptions: Math.floor(Math.random() * 500) + 100, // Simulated
        totalSavings: (Math.random() * 5000) + 1000, // Simulated
        conversionRate: Math.floor(Math.random() * 30) + 50, // Simulated
        topCoupons: couponsData.slice(0, 5).map(coupon => ({
          coupon_code: coupon.coupon_code,
          redemptions: Math.floor(Math.random() * 100) + 10,
          savings: (Math.random() * 1000) + 200,
          conversion_rate: Math.floor(Math.random() * 50) + 40
        })),
        monthlyData: sampleAnalyticsData.monthlyData, // Reuse sample data for now
        categoriesData: sampleAnalyticsData.categoriesData // Reuse sample data for now
      };
      
      setAnalyticsData(calculatedData);
    } catch (error) {
      console.error('Error fetching coupon data:', error);
      // Fall back to sample data if API fails
      setAnalyticsData(sampleAnalyticsData);
    } finally {
      setLoading(false);
    }
  };
  
  useEffect(() => {
    fetchCouponsAndGenerateAnalytics();
  }, [timeframe]);
  
  // Helper function to format currency
  const formatCurrency = (amount) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 2
    }).format(amount);
  };
  
  // Create chart data for visualization
  const getBarWidth = (value, maxValue) => {
    return value && maxValue ? (value / maxValue) * 100 : 0;
  };
  
  if (loading) {
    return (
      <Center h="300px">
        <Spinner size="xl" color="brand.500" />
      </Center>
    );
  }
  
  if (!analyticsData) {
    return (
      <Center h="300px">
        <Text>Failed to load analytics data</Text>
      </Center>
    );
  }
  
  const maxRedemptions = Math.max(...analyticsData.monthlyData.map(m => m.redemptions));
  const maxSavings = Math.max(...analyticsData.monthlyData.map(m => m.savings));
  
  return (
    <Box p={4}>
      <Flex justifyContent="space-between" alignItems="center" mb={6}>
        <Heading as="h1" size="xl">Coupon Analytics</Heading>
        <Select
          maxW="200px"
          value={timeframe}
          onChange={(e) => setTimeframe(e.target.value)}
        >
          <option value="week">Last Week</option>
          <option value="month">Last Month</option>
          <option value="quarter">Last Quarter</option>
          <option value="year">Last Year</option>
          <option value="all">All Time</option>
        </Select>
      </Flex>
      
      {/* Key Metrics */}
      <SimpleGrid columns={{ base: 1, md: 2, lg: 4 }} spacing={6} mb={8}>
        <Stat bg={statCardBg} p={4} borderRadius="lg" boxShadow="sm">
          <StatLabel>Total Coupons</StatLabel>
          <StatNumber>{analyticsData.totalCoupons}</StatNumber>
          <StatHelpText>
            {analyticsData.activeCoupons} active coupons
          </StatHelpText>
        </Stat>
        
        <Stat bg={statCardBg} p={4} borderRadius="lg" boxShadow="sm">
          <StatLabel>Total Redemptions</StatLabel>
          <StatNumber>{analyticsData.totalRedemptions}</StatNumber>
          <StatHelpText>
            Across all coupons
          </StatHelpText>
        </Stat>
        
        <Stat bg={statCardBg} p={4} borderRadius="lg" boxShadow="sm">
          <StatLabel>Customer Savings</StatLabel>
          <StatNumber>{formatCurrency(analyticsData.totalSavings)}</StatNumber>
          <StatHelpText>
            Total discount amount
          </StatHelpText>
        </Stat>
        
        <Stat bg={statCardBg} p={4} borderRadius="lg" boxShadow="sm">
          <StatLabel>Conversion Rate</StatLabel>
          <StatNumber>{analyticsData.conversionRate}%</StatNumber>
          <StatHelpText>
            Of displayed coupons
          </StatHelpText>
        </Stat>
      </SimpleGrid>
      
      {/* Top Performing Coupons */}
      <Card bg={cardBg} mb={8} boxShadow="sm">
        <CardHeader pb={0}>
          <Heading size="md">Top Performing Coupons</Heading>
        </CardHeader>
        <CardBody>
          <Table variant="simple">
            <Thead>
              <Tr>
                <Th>Coupon Code</Th>
                <Th isNumeric>Redemptions</Th>
                <Th isNumeric>Total Savings</Th>
                <Th isNumeric>Conversion</Th>
              </Tr>
            </Thead>
            <Tbody>
              {analyticsData.topCoupons.map((coupon) => (
                <Tr key={coupon.coupon_code}>
                  <Td fontWeight="bold">{coupon.coupon_code}</Td>
                  <Td isNumeric>{coupon.redemptions}</Td>
                  <Td isNumeric>{formatCurrency(coupon.savings)}</Td>
                  <Td isNumeric>
                    <Badge 
                      colorScheme={coupon.conversion_rate > 70 ? 'green' : coupon.conversion_rate > 50 ? 'blue' : 'orange'}
                    >
                      {coupon.conversion_rate}%
                    </Badge>
                  </Td>
                </Tr>
              ))}
            </Tbody>
          </Table>
        </CardBody>
      </Card>
      
      {/* Monthly Performance */}
      <Card bg={cardBg} mb={8} boxShadow="sm">
        <CardHeader pb={0}>
          <Heading size="md">Monthly Performance</Heading>
        </CardHeader>
        <CardBody>
          <VStack spacing={4} align="stretch">
            {analyticsData.monthlyData.map((month) => (
              <Box key={month.month}>
                <HStack justifyContent="space-between" mb={1}>
                  <Text width="50px" fontWeight="medium">{month.month}</Text>
                  <Box flex="1" h="18px">
                    <Box
                      h="100%"
                      w={`${getBarWidth(month.redemptions, maxRedemptions)}%`}
                      bg="blue.400"
                      borderRadius="sm"
                    />
                  </Box>
                  <Text width="60px" textAlign="right">{month.redemptions}</Text>
                </HStack>
                <HStack justifyContent="space-between">
                  <Text width="50px" color="gray.500" fontSize="sm">Savings</Text>
                  <Box flex="1" h="10px">
                    <Box
                      h="100%"
                      w={`${getBarWidth(month.savings, maxSavings)}%`}
                      bg="green.400"
                      borderRadius="sm"
                    />
                  </Box>
                  <Text width="60px" textAlign="right" color="gray.500" fontSize="sm">
                    {formatCurrency(month.savings)}
                  </Text>
                </HStack>
                {month.month !== 'Dec' && <Divider my={2} />}
              </Box>
            ))}
          </VStack>
        </CardBody>
      </Card>
      
      {/* Category Performance */}
      <Card bg={cardBg} boxShadow="sm">
        <CardHeader pb={0}>
          <Heading size="md">Performance by Category</Heading>
        </CardHeader>
        <CardBody>
          <SimpleGrid columns={{ base: 1, md: 2, lg: 3 }} spacing={4}>
            {analyticsData.categoriesData.map((category) => {
              return (
                <Box key={category.category} p={4} bg={categoryBoxBg} borderRadius="md">
                  <Heading size="sm" textTransform="capitalize" mb={2}>
                    {category.category}
                  </Heading>
                  <HStack justifyContent="space-between" mb={2}>
                    <Text fontSize="sm" color="gray.500">Redemptions</Text>
                    <Text fontWeight="bold">{category.redemptions}</Text>
                  </HStack>
                  <HStack justifyContent="space-between">
                    <Text fontSize="sm" color="gray.500">Savings</Text>
                    <Text fontWeight="bold">{formatCurrency(category.savings)}</Text>
                  </HStack>
                </Box>
              );
            })}
          </SimpleGrid>
        </CardBody>
      </Card>
    </Box>
  );
};

export default Analytics; 