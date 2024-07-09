use chrono::{Utc};
use rust_decimal::Decimal;
use rust_decimal::prelude::ToPrimitive;
use crate::entities::{event, event::Entity as Event};
use crate::entities::{event::Model as EventModel};
use sea_orm::*;
use sea_orm::prelude::Expr;
use crate::entities::event::EventType;
use crate::entities::event::EventType::{BUY, UPLOAD, TRANSFER};
use crate::messages::ShishaMessage;

pub struct Query;

impl Query {
    pub async fn insert_message(message: ShishaMessage, db: &DatabaseConnection) -> Result<(), DbErr> {
        let model = event::ActiveModel {
            actor: ActiveValue::Set(message.user),
            event_date: ActiveValue::Set(Utc::now().naive_utc()),
            event_type: ActiveValue::Set(Self::convert_to_type(message.r#type.as_str())),
            target: ActiveValue::Set(message.target),
            amount: ActiveValue::Set(message.amount),
            image: ActiveValue::Set(message.image_uuid),
            ..Default::default()
        };
        Event::insert(model).exec(db).await?;
        Ok(())
    }

    fn convert_to_type(message_type: &str) -> EventType {
        match message_type {
            "upload" => UPLOAD,
            "buy" => BUY,
            "transfer" => TRANSFER,
            t => {
                log::error!("Unknown event type: {}, defaulting to UPLOAD", t);
                UPLOAD
            }
        }
    }
    pub async fn last_10_events(db: &DatabaseConnection) -> Result<Vec<EventModel>, DbErr> {
        Self::last_10_events_query().all(db).await
    }

    pub async fn top_buyer_name(db: &DatabaseConnection) -> Result<Option<String>, DbErr> {
        match Self::top_buyer_query().into_tuple::<(String, i64)>().one(db).await? {
            Some((name, _)) => Ok(Some(name)),
            None => Ok(None)
        }
    }


    pub async fn top_uploader(db: &DatabaseConnection) -> Result<Option<String>, DbErr> {
        match Self::top_uploader_query().into_tuple::<(String, i64)>().one(db).await? {
            Some((name, _)) => Ok(Some(name)),
            None => Ok(None)
        }
    }

    pub async fn count_uploads(db: &DatabaseConnection) -> Result<i64, DbErr> {
        match Self::upload_count_query().into_tuple::<i64>().one(db).await? {
            Some(count) => Ok(count),
            None => Ok(0)
        }
    }
    pub async fn count_money(db: &DatabaseConnection) -> Result<i64, DbErr> {
        match Self::money_earned_query().into_tuple::<Option<Decimal>>().one(db).await? {
            Some(Some(val)) => Ok(val.round().to_i64().expect("Failed to convert number")),
            None | Some(None) => Ok(0i64)
        }
    }
    fn last_10_events_query() -> Select<Event> {
        Event::find().order_by(event::Column::EventDate, Order::Desc).limit(10)
    }

    fn top_buyer_query() -> Select<Event> {
        Self::top_n_by_column(event::Column::Actor, BUY, 1)
    }

    fn top_uploader_query() -> Select<Event> {
        Self::top_n_by_column(event::Column::Actor, UPLOAD, 1)
    }

    fn upload_count_query() -> Select<Event> {
        Event::find()
            .select_only()
            .column_as(event::Column::Id.count(), "events_count")
            .filter(event::Column::EventType.eq(UPLOAD))
    }

    fn money_earned_query() -> Select<Event> {
        Event::find()
            .select_only()
            .column_as(event::Column::Amount.sum(), "total")
            .filter(event::Column::EventType.eq(BUY))
    }

    fn top_n_by_column(col: event::Column, t: EventType, count: u64) -> Select<Event> {
        Event::find()
            .select_only()
            .column(col)
            .column_as(event::Column::Id.count(), "events_count")
            .filter(event::Column::EventType.eq(t))
            .group_by(col)
            .order_by_desc(Expr::cust("events_count"))
            .limit(count)
    }
}

#[cfg(test)]
mod tests {
    use sea_orm::*;
    use crate::service::query::Query;

    #[test]
    fn last_10_events_query_test() {
        let query = "SELECT `event`.`id`, \
        `event`.`actor`, \
        `event`.`event_date`, \
        `event`.`event_type`, \
        `event`.`image`, \
        `event`.`amount`, \
        `event`.`target` \
        FROM `event` \
        ORDER BY `event`.`event_date` DESC LIMIT 10";
        assert_eq!(Query::last_10_events_query().build(DatabaseBackend::MySql).to_string(), query);
    }

    #[test]
    fn top_buyer_query_test() {
        let query = "SELECT \
        `event`.`actor`, \
        COUNT(`event`.`id`) AS `events_count` \
        FROM `event` \
        WHERE `event`.`event_type` = 'BUY' \
        GROUP BY `event`.`actor` \
        ORDER BY events_count DESC \
        LIMIT 1";
        assert_eq!(Query::top_buyer_query().build(DatabaseBackend::MySql).to_string(), query);
    }

    #[test]
    fn upload_count_query_test() {
        let query = "SELECT \
        COUNT(`event`.`id`) AS `events_count` \
        FROM `event` \
        WHERE `event`.`event_type` = 'UPLOAD'";
        assert_eq!(Query::upload_count_query().build(DatabaseBackend::MySql).to_string(), query);
    }

    #[test]
    fn money_earned_query_test() {
        let query = "SELECT \
        SUM(`event`.`amount`) AS `total` \
        FROM `event` \
        WHERE `event`.`event_type` = 'BUY'";
        assert_eq!(Query::money_earned_query().build(DatabaseBackend::MySql).to_string(), query);
    }
}