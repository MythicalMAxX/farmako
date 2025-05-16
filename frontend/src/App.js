import React from 'react';
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom';
import {
  Box,
  Flex,
  HStack,
  Text,
  IconButton,
  Button,
  useDisclosure,
  VStack,
  Heading,
  Container,
  Drawer,
  DrawerBody,
  DrawerHeader,
  DrawerOverlay,
  DrawerContent,
  DrawerCloseButton,
} from '@chakra-ui/react';
import { HamburgerIcon, CloseIcon } from '@chakra-ui/icons';
import CustomerView from './pages/CustomerView';
import AdminDashboard from './pages/AdminDashboard';
import Analytics from './pages/Analytics';

// Create a Home component
const Home = () => (
  <Box p={8} textAlign="center">
    <Heading as="h1" size="xl" mb={6}>
      Welcome to Farmako Coupon System
    </Heading>
    <Text fontSize="lg" mb={6}>
      A comprehensive platform for managing medicine coupons and discounts
    </Text>
    <Button as={Link} to="/customer" colorScheme="brand" size="lg">
      Go to Customer View
    </Button>
  </Box>
);

// Navigation component
const Navigation = () => {
  const { isOpen, onOpen, onClose } = useDisclosure();

  const Links = [
    { name: 'Home', path: '/' },
    { name: 'Customer View', path: '/customer' },
    { name: 'Admin Dashboard', path: '/admin' },
    { name: 'Analytics', path: '/analytics' },
  ];

  const NavLink = ({ children, to }) => (
    <Button
      as={Link}
      px={2}
      py={1}
      rounded={'md'}
      to={to}
      _hover={{
        textDecoration: 'none',
        bg: 'brand.100',
      }}
      variant="ghost"
      onClick={isOpen ? onClose : undefined}
    >
      {children}
    </Button>
  );

  return (
    <Box bg="white" px={4} boxShadow="sm">
      <Flex h={16} alignItems={'center'} justifyContent={'space-between'}>
        <IconButton
          size={'md'}
          icon={isOpen ? <CloseIcon /> : <HamburgerIcon />}
          aria-label={'Open Menu'}
          display={{ md: 'none' }}
          onClick={isOpen ? onClose : onOpen}
        />
        <HStack spacing={8} alignItems={'center'}>
          <Box>
            <Heading size="md" color="brand.500">
              Farmako
            </Heading>
          </Box>
          <HStack as={'nav'} spacing={4} display={{ base: 'none', md: 'flex' }}>
            {Links.map((link) => (
              <NavLink key={link.name} to={link.path}>
                {link.name}
              </NavLink>
            ))}
          </HStack>
        </HStack>
      </Flex>

      {/* Mobile menu */}
      <Drawer
        isOpen={isOpen}
        placement="left"
        onClose={onClose}
        size="xs"
      >
        <DrawerOverlay />
        <DrawerContent>
          <DrawerCloseButton />
          <DrawerHeader borderBottomWidth="1px">Menu</DrawerHeader>
          <DrawerBody>
            <VStack spacing={4} align="stretch">
              {Links.map((link) => (
                <NavLink key={link.name} to={link.path}>
                  {link.name}
                </NavLink>
              ))}
            </VStack>
          </DrawerBody>
        </DrawerContent>
      </Drawer>
    </Box>
  );
};

// Main App component
function App() {
  return (
    <Router>
      <Box minH="100vh">
        <Navigation />
        <Container maxW="container.xl" pt={5}>
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/customer" element={<CustomerView />} />
            <Route path="/admin" element={<AdminDashboard />} />
            <Route path="/analytics" element={<Analytics />} />
          </Routes>
        </Container>
      </Box>
    </Router>
  );
}

export default App; 