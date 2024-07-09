use sea_orm_migration::prelude::*;

pub struct Migration;

impl MigrationName for Migration {
    fn name(&self) -> &str {
        "m_20240624_000001_create_events_table"
    }
}

#[async_trait::async_trait]
impl MigrationTrait for Migration {
    async fn up(&self, manager: &SchemaManager) -> Result<(), DbErr> {
        manager.create_table(
            Table::create().table(Event::Table)
                .col(ColumnDef::new(Event::Id).big_integer().not_null().auto_increment().primary_key())
                .col(ColumnDef::new(Event::Actor).string_len(255).not_null())
                .col(ColumnDef::new(Event::EventDate).date_time().not_null())
                .col(ColumnDef::new(Event::EventType).string_len(255).not_null())
                .col(ColumnDef::new(Event::Image).string_len(255).null())
                .col(ColumnDef::new(Event::Amount).big_integer().null())
                .col(ColumnDef::new(Event::Target).string_len(255).null())
                
                .to_owned()).await?;
        manager.create_index(Index::create().name("IDX_EVENT_EVENT_DATE").table(Event::Table).col(Event::EventDate).to_owned()).await
    }

    async fn down(&self, _manager: &SchemaManager) -> Result<(), DbErr> {
        _manager.drop_index(Index::drop().name("IDX_EVENT_EVENT_DATE").table(Event::Table).to_owned()).await?;
        _manager.drop_table(Table::drop().table(Event::Table).to_owned()).await
    }
}

#[derive(Iden)]
pub enum Event {
    Table,
    Id,
    Actor,
    EventDate,
    EventType,
    Image,
    Amount,
    Target
}