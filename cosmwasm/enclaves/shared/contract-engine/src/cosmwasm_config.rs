/// Right now ContractOperation is used to detect queris and prevent state changes
#[derive(Clone, Copy, Debug)]
pub enum ContractOperation {
    Init,
    Handle,
    Query,
}

#[allow(unused)]
impl ContractOperation {
    pub fn is_init(&self) -> bool {
        matches!(self, ContractOperation::Init)
    }

    pub fn is_handle(&self) -> bool {
        matches!(self, ContractOperation::Handle)
    }

    pub fn is_query(&self) -> bool {
        matches!(self, ContractOperation::Query)
    }
}

//pub const MAX_LOG_LENGTH: usize = 8192;