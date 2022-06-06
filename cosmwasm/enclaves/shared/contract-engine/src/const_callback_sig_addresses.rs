
//DISTR ACCOUNT MODDULE ADDRESS IS HARDCODED
#[cfg(feature = "production")]
pub const COMMUNITY_POOL_ADDR: &str = "trust1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8wz6fau";
#[cfg(not(feature = "production"))]
pub const COMMUNITY_POOL_ADDR: &str = "trust1jv65s3grqf6v6jl3dp4t6c9t9rk99cd8wz6fau";
