package keeper

import (
    "fmt"
    "github.com/NetCloth/netcloth-chain/codec"
    "github.com/NetCloth/netcloth-chain/modules/aipal/types"
    "github.com/NetCloth/netcloth-chain/modules/params"
    sdk "github.com/NetCloth/netcloth-chain/types"
)

type Keeper struct {
    storeKey        sdk.StoreKey
    cdc             *codec.Codec
    supplyKeeper    types.SupplyKeeper
    paramstore      params.Subspace
    codespace       sdk.CodespaceType
}

func NewKeeper(storeKey sdk.StoreKey, cdc *codec.Codec, supplyKeeper types.SupplyKeeper, paramstore params.Subspace, codespace sdk.CodespaceType) Keeper {
    if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
        panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
    }

    return Keeper {
        storeKey:     storeKey,
        cdc:          cdc,
        supplyKeeper: supplyKeeper,
        paramstore:   paramstore.WithKeyTable(ParamKeyTable()),
        codespace:    codespace,
    }
}

func (k Keeper) GetServiceNode(ctx sdk.Context, operator sdk.AccAddress) (obj types.ServiceNode, found bool) {
    store := ctx.KVStore(k.storeKey)
    value := store.Get(types.GetServiceNodeKey(operator))
    if value == nil {
        return obj, false
    }

    obj = types.MustUnmarshalServerNodeObject(k.cdc, value)
    return obj, true
}

func (k Keeper) setServiceNode(ctx sdk.Context, obj types.ServiceNode) {
    store := ctx.KVStore(k.storeKey)
    bz := types.MustMarshalServerNodeObject(k.cdc, obj)
    store.Set(types.GetServiceNodeKey(obj.OperatorAddress), bz)
}

func (k Keeper) delServiceNode(ctx sdk.Context, accAddress sdk.AccAddress) {
    store := ctx.KVStore(k.storeKey)
    store.Delete(types.GetServiceNodeKey(accAddress))
}

func (k Keeper) setServiceNodeByBond(ctx sdk.Context, obj types.ServiceNode) {
    store := ctx.KVStore(k.storeKey)
    bz := types.MustMarshalServerNodeObject(k.cdc, obj)
    store.Set(types.GetServiceNodeByBondKey(obj), bz)
}

func (k Keeper) delServiceNodeByBond(ctx sdk.Context, obj types.ServiceNode) {
    store := ctx.KVStore(k.storeKey)
    store.Delete(types.GetServiceNodeByBondKey(obj))
}

func (k Keeper) setServiceNodeByMonikerIndex(ctx sdk.Context, obj types.ServiceNode) {
    store := ctx.KVStore(k.storeKey)
    store.Set(types.GetServiceNodeByMonikerKey(obj.Moniker), obj.OperatorAddress)
}

func (k Keeper) getServiceNodeAddByMoniker(ctx sdk.Context, moniker string) (acc sdk.AccAddress, exist bool) {
    store := ctx.KVStore(k.storeKey)
    v := store.Get(types.GetServiceNodeByMonikerKey(moniker))
    return v, v != nil
}

func (k Keeper) delServiceNodeByMonikerIndex(ctx sdk.Context, moniker string) {
    store := ctx.KVStore(k.storeKey)
    store.Delete(types.GetServiceNodeByMonikerKey(moniker))
}

func (k Keeper) createServiceNode(ctx sdk.Context, m types.MsgServiceNodeClaim) {
    n := types.NewServiceNode(m.OperatorAddress, m.Moniker, m.Website, types.ServiceType(m.ServiceType), m.ServerEndPoint, m.Details, m.Bond)
    k.setServiceNode(ctx, n)
    k.setServiceNodeByBond(ctx, n)
    k.setServiceNodeByMonikerIndex(ctx, n)
}

func (k Keeper) updateServiceNode(ctx sdk.Context, old types.ServiceNode, new types.MsgServiceNodeClaim) {
    u := types.NewServiceNode(new.OperatorAddress, new.Moniker, new.Website, types.ServiceType(new.ServiceType), new.ServerEndPoint, new.Details, new.Bond)
    k.setServiceNode(ctx, u)

    k.delServiceNodeByBond(ctx, old)
    k.setServiceNodeByBond(ctx, u)

    k.delServiceNodeByMonikerIndex(ctx, old.Moniker)
    k.setServiceNodeByMonikerIndex(ctx, u)
}

func (k Keeper) deleteServiceNode(ctx sdk.Context, n types.ServiceNode) {
    k.delServiceNode(ctx, n.OperatorAddress)
    k.delServiceNodeByBond(ctx, n)
    k.delServiceNodeByMonikerIndex(ctx, n.Moniker)
}

func (k Keeper) bond(ctx sdk.Context, aa sdk.AccAddress, amt sdk.Coin) sdk.Error {
    return k.supplyKeeper.SendCoinsFromAccountToModule(ctx, aa, types.ModuleName, sdk.Coins{amt})
}


/*
founded {
    bond >= minBond {
        bond > currentBond {
            ensure Bond(bond - currentBond)
        } else (bond < currentBond) {
            UnBond(currentBond - bond)
        } else {
        }
        updateServiceNode
    } else {
        UnBond
        deleteServiceNode
    }
} else {
    bond >= minBond {
        ensure moniker uniq
        ensure Bond
        createServiceNode
    } else {
        return err
    }
}
*/
func (k Keeper) DoServiceNodeClaim(ctx sdk.Context, m types.MsgServiceNodeClaim) (err sdk.Error) {
    acc, monikerExist := k.getServiceNodeAddByMoniker(ctx, m.Moniker)
    if monikerExist && !acc.Equals(m.OperatorAddress) {
        return types.ErrMonikerExist(fmt.Sprintf("moniker: [%s] already exist", m.Moniker))
    }

    minBond := k.GetMinBond(ctx)
    n, found := k.GetServiceNode(ctx, m.OperatorAddress)
    if found {
        if m.Bond.IsGTE(minBond) {
            if n.Bond.IsLT(m.Bond) {
                err := k.bond(ctx, m.OperatorAddress, m.Bond.Sub(n.Bond))
                if err != nil {
                    return err
                }
            } else if m.Bond.IsLT(n.Bond) {
                k.toUnbondingQueue(ctx, m.OperatorAddress, n.Bond.Sub(m.Bond))
            } else {
            }
            k.updateServiceNode(ctx, n, m)
        } else {
            k.toUnbondingQueue(ctx, m.OperatorAddress, n.Bond)
            k.deleteServiceNode(ctx, n)
        }
    } else {
        if m.Bond.IsGTE(minBond) {
            err := k.bond(ctx, m.OperatorAddress, m.Bond)
            if err != nil  {
                return err
            }

            k.createServiceNode(ctx, m)
        } else {
            return types.ErrBondInsufficient(fmt.Sprintf("bond insufficient, min bond: %s, actual bond: %s", minBond.String(), m.Bond.String()))
        }
    }

    return nil
}

func (k Keeper) GetAllServerNodes(ctx sdk.Context) (serverNodes types.ServiceNodes) {
    store := ctx.KVStore(k.storeKey)
    iterator := sdk.KVStorePrefixIterator(store, types.ServiceNodeByBondKey)
    defer iterator.Close()

    for ; iterator.Valid(); iterator.Next() {
        validator := types.MustUnmarshalServerNodeObject(k.cdc, iterator.Value())
        serverNodes = append(serverNodes, validator)
    }
    return serverNodes
}