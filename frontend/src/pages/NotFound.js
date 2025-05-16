import React from 'react';
import { Box, Heading, Text, Button, VStack, useColorModeValue } from '@chakra-ui/react';
import { Link as RouterLink } from 'react-router-dom';

const NotFound = () => {
  const bgColor = useColorModeValue('white', 'gray.800');
  const borderColor = useColorModeValue('gray.200', 'gray.700');

  return (
    <Box 
      p={8} 
      maxW="md" 
      mx="auto" 
      mt={20} 
      borderWidth="1px" 
      borderRadius="lg" 
      boxShadow="lg" 
      bg={bgColor}
      borderColor={borderColor}
    >
      <VStack spacing={6} align="center">
        <Heading as="h1" size="2xl">404</Heading>
        <Heading as="h2" size="md">Page Not Found</Heading>
        <Text align="center">
          The page you're looking for doesn't exist or has been moved.
        </Text>
        <Button 
          as={RouterLink} 
          to="/" 
          colorScheme="brand"
          size="lg"
        >
          Back to Home
        </Button>
      </VStack>
    </Box>
  );
};

export default NotFound; 