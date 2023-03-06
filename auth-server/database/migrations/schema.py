from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy import (
    Column,
    String,
    Sequence,
    Integer,
)

Base = declarative_base()

class Auth(Base):
    __tablename__ = 'authentications'
    id = Column(Integer, Sequence("auth_id"), primary_key=True)
    token = Column(String, nullable=False, unique=True)